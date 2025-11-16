package consumers

import (
	"context"
	"git.abanppc.com/farin-project/event/config"
	"git.abanppc.com/farin-project/event/domain/biz/event"
	v1 "git.abanppc.com/farin-project/event/domain/entity/v1"
	"git.abanppc.com/farin-project/event/pkg/rabbit"
	"github.com/rabbitmq/amqp091-go"
	"log/slog"
	"strings"
	"time"
)

const workersCount = 1
const msgsQueueLength = 5

type Event struct {
	rbt       *rabbit.Rabbit
	delivery  <-chan amqp091.Delivery
	logger    *slog.Logger
	msgsQueue chan amqp091.Delivery
	rbtConfig *config.RabbitMQ
	eventBiz  event.EventBizInterface
}

func NewEvent(l *slog.Logger, rbt *rabbit.Rabbit, rbtConfig *config.RabbitMQ, eventBiz event.EventBizInterface) *Event {
	return &Event{
		rbt:       rbt,
		logger:    l,
		msgsQueue: make(chan amqp091.Delivery, msgsQueueLength),
		rbtConfig: rbtConfig,
		eventBiz:  eventBiz,
	}
}

func (u *Event) Setup() error {
	d, err := u.rbt.Ch.Consume(u.rbtConfig.EventQueue, "", false,
		false, false, false, nil)
	if err != nil {
		return err
	}
	u.delivery = d
	return nil
}

func (u *Event) Consume() {
	for i := 0; i < workersCount; i++ {
		go u.worker()
	}
	for msg := range u.delivery {
		u.msgsQueue <- msg
	}
}

func (u *Event) worker() {
	lg := u.logger.With("method", "worker", "queue", u.rbtConfig.EventQueue)
	lg.Info("worker started")
	for msg := range u.msgsQueue {
		u.logger.Info("rabbit message received in msg queue go channel: routingKey:%s,body", msg.RoutingKey)
		routingKey := msg.RoutingKey
		// Check if the routing key matches a known dynamic pattern
		if strings.HasPrefix(routingKey, "farin.vehicles.") && strings.HasSuffix(routingKey, ".event") {
			// Split the routing key to extract dynamic parts
			parts := strings.Split(routingKey, ".")
			if len(parts) == 6 {
				lprVehicleId := parts[2]
				userId := parts[4]
				lg.Info("Processing for event creation", "vehicle id", lprVehicleId, "user id", userId)
				ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
				ev := v1.Event{}
				if err := ev.UnmarshalNonFileFieldsJSON(msg.Body); err != nil {
					lg.Error("failed to unmarshal event", "error", err)
					return
				}
				if err := u.eventBiz.CreateEvent(ctx, &ev); err != nil {
					Nack(msg, lg)
					continue
				}
				if err := msg.Ack(false); err != nil {
					lg.Error("failed to ack message", "error", err)
				}
				lg.Info("Acked event creation", "vehicle id", lprVehicleId, "user id", userId)
				continue
			}
		}
		Nack(msg, lg)
		lg.Warn("rabbit message nacked", "routingKey", msg.RoutingKey)
	}
}

func Nack(msg amqp091.Delivery, lg *slog.Logger) {
	if err := msg.Nack(false, true); err != nil {
		lg.Error("failed to nack message", "error", err, "routingKey", msg.RoutingKey)
	}
	time.Sleep(10 * time.Millisecond) // cool down
}
