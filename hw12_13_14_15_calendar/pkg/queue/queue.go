package queue

import (
	"context"
	"time"
)

type Config struct {
	QueueName      string        `mapstructure:"queue_name"`
	ExchangeName   string        `mapstructure:"exchange_name"`
	User           string        `mapstructure:"user"`
	Password       string        `mapstructure:"password"`
	Host           string        `mapstructure:"host"`
	Port           int           `mapstructure:"port"`
	ReconnectDelay time.Duration `mapstructure:"reconnect_delay"`
	ConsumerName   string        `mapstructure:"consumer_name"`
	MaxRetries     uint64        `mapstructure:"max_retries"`
}

type Message struct {
	Body   []byte
	Ack    func(multiple bool) error
	Reject func(requeue bool) error
}
type Connector interface {
	Connect(ctx context.Context) error
}

type Producer interface {
	Connector
	Publish(data []byte) error
}

type Handler interface {
	Handle(ctx context.Context, d *Message) error
}

type Consumer interface {
	Listen(ctx context.Context, h Handler) error
}
