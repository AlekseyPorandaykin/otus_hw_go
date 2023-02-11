package ampq

import (
	"context"
	"time"

	"github.com/AlekseyPorandaykin/otus_hw_go/hw12_13_14_15_calendar/pkg/logger"
	"github.com/AlekseyPorandaykin/otus_hw_go/hw12_13_14_15_calendar/pkg/queue"
	"github.com/streadway/amqp"
)

type Producer struct {
	con *Connection
	log logger.Logger
}

func NewProducer(con *Connection, log logger.Logger) *Producer {
	return &Producer{
		log: log,
		con: con,
	}
}

func (p *Producer) Connect(ctx context.Context) error {
	return p.con.Connect(ctx)
}

func (p *Producer) Publish(data []byte) error {
	if p.con.status == notInitStatus {
		return queue.NotConnectError
	}
	if p.con.status == closeStatus {
		return queue.ClosedError
	}

	return p.con.ch.Publish(
		p.con.cfg.ExchangeName,
		p.con.cfg.QueueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        data,
			Timestamp:   time.Now(),
		})
}
