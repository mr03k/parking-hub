package rabbit

import (
	"errors"
	"git.abanppc.com/farin-project/event/config"
	"github.com/rabbitmq/amqp091-go"
)

type Rabbit struct {
	Ch        *amqp091.Channel
	Conn      *amqp091.Connection
	rbtConfig *config.RabbitMQ
}

func NewRabbit(rbtConfig *config.RabbitMQ) *Rabbit {
	return &Rabbit{rbtConfig: rbtConfig}
}

func (r *Rabbit) Setup() error {
	if r.Conn != nil {
		r.Conn.Close()
	}
	if r.Ch != nil {
		r.Ch.Close()
	}
	conn, err := amqp091.Dial(r.rbtConfig.Host)
	if err != nil {
		return err
	}
	r.Conn = conn
	r.Ch, err = conn.Channel()
	if err != nil {
		return err
	}
	return nil
}

// HealthCheck checks if RabbitMQ is ready by performing a basic ping-like operation
func (r *Rabbit) HealthCheck() error {
	// Check if the connection is open
	if r.Conn == nil || r.Conn.IsClosed() {
		return errors.New("rabbitmq connection is not open")
	}

	// Check if the channel is valid
	if r.Ch == nil || r.Ch.IsClosed() {
		return errors.New("rabbitmq channel is not open")
	}
	return nil
}
