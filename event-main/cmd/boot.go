package main

import (
	"git.abanppc.com/farin-project/event/config"
	"git.abanppc.com/farin-project/event/internal/rabbit/consumers"
	"git.abanppc.com/farin-project/event/pkg/rabbit"

	"log"
)

type Bootstrap struct {
	consumers []consumers.Consumer
	rbtConfig *config.RabbitMQ
}

func NewBootstrap(rbtConfig *config.RabbitMQ, consumers ...consumers.Consumer) *Bootstrap {
	return &Bootstrap{
		rbtConfig: rbtConfig,
		consumers: consumers,
	}
}

func (b *Bootstrap) Boot() {
	rbt := rabbit.NewRabbit(b.rbtConfig)
	if err := rbt.Setup(); err != nil {
		log.Fatal(err.Error())
	}

	for _, consumer := range b.consumers {
		if err := consumer.Setup(); err != nil {
			log.Fatal("setup consumer error: " + err.Error())
		}
		go consumer.Consume()
	}
}
