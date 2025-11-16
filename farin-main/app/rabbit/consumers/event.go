package consumers

import (
	"context"
	"farin/infrastructure/godotenv"
	"github.com/rabbitmq/amqp091-go"
	"log/slog"
	"time"
)

const queueLength = 10
const workersCount = 3

type EventConsumer struct {
	q         *amqp091.Queue
	delivery  <-chan amqp091.Delivery
	logger    *slog.Logger
	msgsQueue chan amqp091.Delivery
	handlers  map[string]HandlerFunc
	env       *godotenv.Env
}

func NewEventConsumer(l *slog.Logger, env *godotenv.Env) *EventConsumer {
	return &EventConsumer{logger: l.With("layer", "EventConsumer"), handlers: map[string]HandlerFunc{},
		msgsQueue: make(chan amqp091.Delivery, queueLength), env: env}
}

func (c *EventConsumer) RunInnerWorkers() {
	for i := 0; i < workersCount; i++ {
		go c.innerWorker()
	}
}

func (c *EventConsumer) Setup(ch *amqp091.Channel) error {
	q, err := ch.QueueDeclare(c.env.RabbitMQGeodataQueue, true, false,
		false, false, map[string]interface{}{
			"x-queue-type": "quorum",
		})
	if err != nil {
		return err
	}
	if err := ch.QueueBind(q.Name, "farin.vehicles.drivers.event", c.env.RabbitMQEventExchange, false, nil); err != nil {
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

func (c *EventConsumer) RegisterHandler(routingKey string, handler HandlerFunc) {
	c.handlers[routingKey] = handler
}

func (c *EventConsumer) Worker() {
	for msg := range c.delivery {
		c.msgsQueue <- msg
	}
}

func (c *EventConsumer) innerWorker() {
	lg := c.logger.With("method", "worker")
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
			ctx, _ := context.WithTimeout(context.Background(), time.Second*65)
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
