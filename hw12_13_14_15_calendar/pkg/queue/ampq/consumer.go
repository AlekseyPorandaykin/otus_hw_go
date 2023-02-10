package ampq

import (
	"context"

	"github.com/AlekseyPorandaykin/otus_hw_go/hw12_13_14_15_calendar/pkg/logger"
	"github.com/AlekseyPorandaykin/otus_hw_go/hw12_13_14_15_calendar/pkg/queue"
	"go.uber.org/zap"
)

type Consumer struct {
	con              *Connection
	log              logger.Logger
	createConsumerCh chan interface{}
	closeConsumerCh  chan interface{}
}

func NewConsumer(con *Connection, log logger.Logger) *Consumer {
	return &Consumer{
		con:              con,
		log:              log,
		createConsumerCh: make(chan interface{}),
		closeConsumerCh:  make(chan interface{}),
	}
}

func (c *Consumer) Listen(ctx context.Context, h queue.Handler) error {
	if err := c.con.Connect(ctx); err != nil {
		return queue.ConsumeError.Wrap(err)
	}
	defer c.con.Close()
	c.con.AddCreateSubscriber(c.createConsumerCh)
	c.con.AddCloseSubscriber(c.closeConsumerCh)
	deliveries, errC := c.consume(ctx)
	for {
		select {
		case err := <-errC:
			return queue.ConsumeError.Wrap(err)
		case <-ctx.Done():
			return nil
		case message := <-deliveries:
			if message.Body == nil {
				continue
			}
			if errH := h.Handle(message); errH != nil {
				c.log.Error("Error handle message", zap.Error(errH))
				if errR := message.Reject(false); errR != nil {
					return queue.ConsumeError.Wrap(errR)
				}
				continue
			}
			if errA := message.Ack(false); errA != nil {
				return queue.ConsumeError.Wrap(errA)
			}
		}
	}
}

func (c *Consumer) consume(ctx context.Context) (<-chan *queue.Message, chan error) {
	errCh := make(chan error)
	messageCh := make(chan *queue.Message)
	callback := c.makeConsumerCallback(errCh, messageCh)

	go func(consumer func(ctx context.Context)) {
		go consumer(ctx)
		for {
			select {
			case <-ctx.Done():
				return
			case <-c.createConsumerCh:
				go consumer(ctx)
			default:
				if c.con.status == closeStatus {
					errCh <- queue.ClosedError
					return
				}
			}
		}
	}(callback)

	return messageCh, errCh
}

func (c *Consumer) makeConsumerCallback(errCh chan error, messageCh chan *queue.Message) func(context.Context) {
	return func(ctx context.Context) {
		c.log.Debug("Connect consumer to broker")
		defer c.log.Debug("Close consumer")
		if c.con.status == notInitStatus {
			errCh <- queue.NotConnectError
			return
		}
		d, errC := c.con.ch.Consume(c.con.cfg.QueueName, c.con.cfg.ConsumerName, false, false, false, false, nil)
		if errC != nil {
			errCh <- queue.ConsumeError.Wrap(errC)
			return
		}
		for {
			select {
			case <-ctx.Done():
				return
			case <-c.closeConsumerCh:
				return
			case delivery := <-d:
				messageCh <- &queue.Message{
					Body:   delivery.Body,
					Ack:    delivery.Ack,
					Reject: delivery.Reject,
				}
			}
		}
	}
}
