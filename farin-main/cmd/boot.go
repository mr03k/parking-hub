package main

import (
	"farin/app/api/routes"
	"farin/app/api/validators"
	"farin/app/rabbit/consumers"
	"farin/app/rabbit/handlers"
	"farin/docs"
	_ "farin/docs"
	"farin/infrastructure/godotenv"
	"farin/infrastructure/rabbit"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
	"time"
)

type Boot struct {
	rts           []routes.Router
	validators    validators.Validators
	eventConsumer *consumers.EventConsumer
	event         *handlers.VehicleRecord
	rbt           *rabbit.Rabbit
	cr            *rabbit.ConsumerRunner
}

func NewBoot(validators validators.Validators, event *handlers.VehicleRecord, eventConsumer *consumers.EventConsumer,
	cr *rabbit.ConsumerRunner, rbt *rabbit.Rabbit, rts ...routes.Router) *Boot {
	return &Boot{rts: rts, validators: validators, event: event, eventConsumer: eventConsumer, cr: cr, rbt: rbt}
}

func (b *Boot) Boot() {
	env := godotenv.NewEnv()
	env.Load()
	r := gin.Default()

	b.validators.Setup()
	for _, router := range b.rts {
		router.SetupRoutes(r)
	}
	docs.SwaggerInfo.BasePath = "/api"

	b.event.RegisterConsumer(b.eventConsumer)
	go func() {
		if err := b.rbt.Setup(b.cr); err != nil {
			log.Fatalf("failed to setup rabbitmq:%s", err)
		}
	}()
	time.Sleep(5 * time.Second) //todo: clean
	b.cr.RunInnerWorkers()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	if err := r.Run(":" + env.HTTPPort); err != nil {
		log.Fatalf("error running gin grpcServer error:%s", err.Error())
	}
}
