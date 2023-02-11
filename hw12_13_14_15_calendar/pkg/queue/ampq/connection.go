package ampq

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/AlekseyPorandaykin/otus_hw_go/hw12_13_14_15_calendar/pkg/logger"
	"github.com/AlekseyPorandaykin/otus_hw_go/hw12_13_14_15_calendar/pkg/queue"
	retry "github.com/sethvargo/go-retry"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

const (
	defaultReconnectDelay      time.Duration = 5 * time.Second
	defaultReconnectMaxRetries uint64        = 10
)

const (
	notInitStatus = iota
	connectedStatus
	closeStatus
)

type Connection struct {
	cfg           *queue.Config
	log           logger.Logger
	con           *amqp.Connection
	ch            *amqp.Channel
	brokerErrChan chan *amqp.Error
	status        int

	muClose          sync.Mutex
	closeSubscribers []chan struct{}

	muCreate          sync.Mutex
	createSubscribers []chan struct{}
	signal            struct{}
}

func (c *Connection) Connect(ctx context.Context) error {
	if err := c.init(); err != nil {
		return queue.NotConnectError.Wrap(err)
	}

	go func() {
		if err := c.reconnect(ctx); err != nil {
			c.log.Error("Error reconnect", zap.Error(err))
			c.status = closeStatus
			return
		}
	}()

	c.setupSubscribes()

	return nil
}

func (c *Connection) AddCreateSubscriber(sub chan struct{}) {
	c.muCreate.Lock()
	defer c.muCreate.Unlock()
	c.createSubscribers = append(c.createSubscribers, sub)
}

func (c *Connection) AddCloseSubscriber(sub chan struct{}) {
	c.muClose.Lock()
	defer c.muClose.Unlock()
	c.closeSubscribers = append(c.closeSubscribers, sub)
}

func (c *Connection) Close() error {
	c.status = closeStatus
	if c.ch != nil {
		c.ch.Close()
	}
	if c.con != nil {
		return c.con.Close()
	}
	return nil
}

func NewConnection(cfg *queue.Config, log logger.Logger) *Connection {
	return &Connection{
		cfg:               cfg,
		log:               log,
		brokerErrChan:     make(chan *amqp.Error),
		createSubscribers: make([]chan struct{}, 0),
		closeSubscribers:  make([]chan struct{}, 0),
		status:            notInitStatus,
		signal:            struct{}{},
	}
}

func (c *Connection) reconnect(ctx context.Context) error {
	retryCallback := c.makeReconnectRetryCallback(ctx)

	d := defaultReconnectDelay
	if c.cfg.ReconnectDelay > 0 {
		d = c.cfg.ReconnectDelay
	}

	ticker := time.NewTicker(d)
	for {
		select {
		case <-ctx.Done():
			return nil
		case err := <-c.brokerErrChan:
			if err != nil {
				c.notifyCloseSubscribers()
				if errR := retryCallback(); errR != nil {
					return errR
				}
				c.notifyCreateSubscribers()
			}
		case <-ticker.C:
			if c.con.IsClosed() && c.status != closeStatus {
				c.notifyCloseSubscribers()
				if errR := retryCallback(); errR != nil {
					return errR
				}
				c.notifyCreateSubscribers()
			}
		}
	}
}

func (c *Connection) makeReconnectRetryCallback(ctx context.Context) func() error {
	maxRetries := defaultReconnectMaxRetries
	if c.cfg.MaxRetries > 0 {
		maxRetries = c.cfg.MaxRetries
	}
	r := retry.NewFibonacci(time.Second * 1)
	r = retry.WithMaxRetries(maxRetries, r)
	r = retry.WithMaxDuration(time.Second*30, r)

	return func() error {
		return retry.Do(ctx, r, func(ctx context.Context) error {
			if errS := c.init(); errS != nil {
				c.log.Debug("Error retry connect", zap.Error(errS))
				return retry.RetryableError(errS)
			}
			c.log.Debug("Reconnect to broker")
			return nil
		})
	}
}

func (c *Connection) init() error {
	if errCon := c.setupConnect(); errCon != nil {
		return errCon
	}
	if errCh := c.setupChannel(); errCh != nil {
		return errCh
	}
	return nil
}

func (c *Connection) setupConnect() error {
	con, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%d/", c.cfg.User, c.cfg.Password, c.cfg.Host, c.cfg.Port))
	if err != nil {
		return queue.NotConnectError.Wrap(err)
	}
	c.con = con
	c.status = connectedStatus

	return nil
}

func (c *Connection) setupChannel() error {
	ch, errC := c.con.Channel()
	if errC != nil {
		return queue.NotOpenChannelError.Wrap(errC)
	}
	if c.cfg.ExchangeName != "" {
		errE := ch.ExchangeDeclare(
			c.cfg.ExchangeName,
			"topic",
			true,
			false,
			false,
			false,
			nil,
		)
		if errE != nil {
			return errE
		}
	}

	if err := ch.Confirm(false); err != nil {
		return err
	}

	_, errQ := ch.QueueDeclare(c.cfg.QueueName, false, false, false, false, nil)
	if errQ != nil {
		return errQ
	}
	if c.cfg.ExchangeName != "" {
		if errB := ch.QueueBind(c.cfg.QueueName, "#", c.cfg.ExchangeName, false, nil); errB != nil {
			return errB
		}
	}
	c.ch = ch

	return nil
}

func (c *Connection) setupSubscribes() {
	c.con.NotifyClose(c.brokerErrChan)
}

func (c *Connection) notifyCreateSubscribers() {
	for _, subscriber := range c.createSubscribers {
		subscriber <- c.signal
	}
}

func (c *Connection) notifyCloseSubscribers() {
	for _, subscriber := range c.closeSubscribers {
		subscriber <- c.signal
	}
}
