package consumers

import (
	"context"
	"git.abanppc.com/farin-project/state/infrastructure/godotenv"
	"github.com/rabbitmq/amqp091-go"
	"log/slog"
	"time"
)

const queueLength = 30
const workersCount = 10

type StateConsumer struct {
	q         *amqp091.Queue
	delivery  <-chan amqp091.Delivery
	logger    *slog.Logger
	handlers  map[string]HandlerFunc
	env       *godotenv.Env
	msgsQueue chan amqp091.Delivery
}

func NewEventConsumer(l *slog.Logger, env *godotenv.Env) *StateConsumer {
	ec := &StateConsumer{logger: l.With("layer", "StateConsumer"), handlers: map[string]HandlerFunc{},
		env: env, msgsQueue: make(chan amqp091.Delivery, queueLength)}
	return ec
}

func (c *StateConsumer) RunInnerWorkers() {
	for i := 0; i < workersCount; i++ {
		go c.innerWorker()
	}
}

func (c *StateConsumer) Setup(ch *amqp091.Channel) error {
	q, err := ch.QueueDeclare(c.env.RabbitMQStateQueue, true, false,
		false, false, map[string]interface{}{
			"x-max-length": 1000,
			"x-overflow":   "reject-publish",
		})
	if err != nil {
		return err
	}
	if err := ch.QueueBind(q.Name, "farin.vehicles.drivers.state", c.env.RabbitMQEventExchange, false, nil); err != nil {
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

func (c *StateConsumer) RegisterHandler(routingKey string, handler HandlerFunc) {
	c.handlers[routingKey] = handler
}

func (c *StateConsumer) Worker() {
	for msg := range c.delivery {
		c.msgsQueue <- msg
	}
}

func (c *StateConsumer) innerWorker() {
	lg := c.logger.With("method", "Worker")
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
			ctx, _ := context.WithTimeout(context.Background(), time.Second*55)
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
