package consumers

import (
	"context"
	"github.com/rabbitmq/amqp091-go"
)

type HandlerFunc func(context.Context, []byte) error

type Consumer interface {
	Setup(ch *amqp091.Channel, rch *amqp091.Channel) error
	Worker()
	RegisterHandler(routingKey string, handler HandlerFunc)
	RunInnerWorkers()
}
