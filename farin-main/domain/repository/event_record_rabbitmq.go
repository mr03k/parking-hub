package repository

import (
	"context"
	"encoding/json"
	"farin/infrastructure/godotenv"
	"farin/infrastructure/rabbit"
	"git.abanppc.com/farin-project/vehicle-records/domain/entity"
	"github.com/rabbitmq/amqp091-go"
)

type EventRecordRabbitMQ struct {
	rbt *rabbit.Rabbit
	env *godotenv.Env
}

func NewEventRecordRabbitMQ(rbt *rabbit.Rabbit, env *godotenv.Env) *EventRecordRabbitMQ {
	return &EventRecordRabbitMQ{rbt: rbt, env: env}
}

func (er EventRecordRabbitMQ) Store(ctx context.Context, record *entity.VehicleRecord) error {
	b, err := json.Marshal(&record)
	if err != nil {
		return err
	}
	if err := er.rbt.Ch.PublishWithContext(ctx, er.env.RabbitMQInternalExchange, "farin.vehicles.drivers.event",
		false, false, amqp091.Publishing{
			ContentType: "application/json",
			Body:        b,
		}); err != nil {
		return err
	}
	return nil
}
