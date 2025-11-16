package main

import (
	"git.abanppc.com/farin-project/state/app/api/routes"
	"git.abanppc.com/farin-project/state/app/rabbit/consumers"
	"git.abanppc.com/farin-project/state/app/rabbit/handlers"
	"git.abanppc.com/farin-project/state/infrastructure/godotenv"
	"git.abanppc.com/farin-project/state/infrastructure/rabbit"
	"github.com/gin-gonic/gin"
	"log"
)

type Boot struct {
	rts           []routes.Router
	stateConsumer *consumers.StateConsumer
	stateHandler  *handlers.State
	rbt           *rabbit.Rabbit
	cr            *rabbit.ConsumerRunner
}

func NewBoot(event *handlers.State, eventConsumer *consumers.StateConsumer, rbt *rabbit.Rabbit,
	cr *rabbit.ConsumerRunner, rts ...routes.Router) *Boot {
	return &Boot{stateHandler: event, stateConsumer: eventConsumer,
		rts: rts, rbt: rbt, cr: cr}
}

func (b *Boot) Boot() {
	env := godotenv.NewEnv()
	env.Load()
	r := gin.Default()
	for _, router := range b.rts {
		router.SetupRoutes(r)
	}

	b.stateHandler.RegisterConsumer(b.stateConsumer)
	b.cr.RunInnerWorkers()
	go func() {
		if err := b.rbt.Setup(b.cr); err != nil {
			log.Fatalf("failed to setup rabbitmq:%s", err)
		}
	}()

	if err := r.Run(":" + env.HTTPPort); err != nil {
		log.Fatalf("error running gin grpcServer error:%s", err.Error())
	}
}
