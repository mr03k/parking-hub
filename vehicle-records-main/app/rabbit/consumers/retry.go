package consumers

import (
	"context"
	"log/slog"

	"git.abanppc.com/farin-project/vehicle-records/infrastructure/godotenv"
	"github.com/rabbitmq/amqp091-go"
)

const retryQueueLength = 8
const retryWorkersCount = 3
const retryQueueTTL = 60000 * 60 * 18 //18 hours

type RetryConsumer struct {
	q         *amqp091.Queue
	delivery  <-chan amqp091.Delivery
	logger    *slog.Logger
	handlers  map[string]HandlerFunc
	env       *godotenv.Env
	msgsQueue chan amqp091.Delivery
}

func NewRetryConsumer(l *slog.Logger, env *godotenv.Env) *RetryConsumer {
	ec := &RetryConsumer{logger: l.With("layer", "RetryConsumer"), handlers: map[string]HandlerFunc{},
		env: env, msgsQueue: make(chan amqp091.Delivery, retryQueueLength)}
	return ec
}

func (c *RetryConsumer) RunInnerWorkers() {
	for i := 0; i < retryWorkersCount; i++ {
		go c.innerWorker()
	}
}

func (c *RetryConsumer) Setup(_ *amqp091.Channel, ch *amqp091.Channel) error {
	q, err := ch.QueueDeclare(c.env.RabbitMQRetryQueue, true, false,
		false, false, map[string]interface{}{
			"x-queue-type":  "quorum",
			"x-message-ttl": retryQueueTTL,
		})
	if err != nil {
		return err
	}
	if err := ch.QueueBind(q.Name, "farin.vehicles.drivers.event.retry", c.env.RabbitMQInternalExchange,
		false, nil); err != nil {
		return err
	}
	c.q = &q
	d, err := ch.Consume(c.q.Name, "", false,
		false, false, false, nil)
	if err != nil {
		return err
	}
	c.delivery = d

	return nil
}

func (c *RetryConsumer) RegisterHandler(routingKey string, handler HandlerFunc) {
	c.handlers[routingKey] = handler
}

func (c *RetryConsumer) Worker() {
	lg := c.logger.With("method", "Worker")
	lg.Info("started RetryConsumer worker")
	for msg := range c.delivery {
		c.msgsQueue <- msg
	}
}

func (c *RetryConsumer) innerWorker() {
	lg := c.logger.With("method", "Worker")
	lg.Info("started RetryConsumer inner worker")

	for msg := range c.msgsQueue {
		lg.Info("rabbit message received in msg queue go channel", slog.String("routingKey", msg.RoutingKey))
		handler, ok := c.handlers[msg.RoutingKey]
		if !ok {
			lg.Warn("no handler found for routingKey", slog.String("routingKey", msg.RoutingKey))
			if err := msg.Ack(false); err != nil {
				lg.Error("failed to ack message", slog.Any("error", err))
			}
			lg.Warn("rabbit message acked(no handler found)", slog.String("routingKey", msg.RoutingKey))
			continue
		} else {
			ctx := context.Background()
			if err := handler(ctx, msg.Body); err == nil {
				if err := msg.Ack(false); err != nil {
					lg.Error("failed to ack message", slog.Any("error", err))
				}
				lg.Info("rabbit message acked", slog.String("routingKey", msg.RoutingKey))
				continue
			}
		}

		if err := msg.Nack(false, true); err != nil {
			lg.Error("failed to nack message", slog.Any("error", err))
			continue
		}
		lg.Warn("rabbit message nacked", slog.String("routingKey", msg.RoutingKey))
	}
}
