package main

import (
	"context"
	"git.abanppc.com/farin-project/vehicle-records/app/api/routes"
	"git.abanppc.com/farin-project/vehicle-records/app/rabbit/consumers"
	"git.abanppc.com/farin-project/vehicle-records/app/rabbit/handlers"
	"git.abanppc.com/farin-project/vehicle-records/cmd/manual_retry"
	docs "git.abanppc.com/farin-project/vehicle-records/docs"
	"git.abanppc.com/farin-project/vehicle-records/domain/service"
	"git.abanppc.com/farin-project/vehicle-records/infrastructure/godotenv"
	"git.abanppc.com/farin-project/vehicle-records/infrastructure/rabbit"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
	"log/slog"
	"time"
)

type Boot struct {
	rts           []routes.Router
	eventConsumer *consumers.EventConsumer
	event         *handlers.EventVehicleRecord
	rbt           *rabbit.Rabbit
	cr            *rabbit.ConsumerRunner
	retryConsumer *consumers.RetryConsumer
	retry         *handlers.Retry
	vs            *service.VehicleRecordService
	lg            *slog.Logger
	env           *godotenv.Env
}

func NewBoot(event *handlers.EventVehicleRecord, eventConsumer *consumers.EventConsumer, rbt *rabbit.Rabbit,
	cr *rabbit.ConsumerRunner, retryConsumer *consumers.RetryConsumer, retry *handlers.Retry, vs *service.VehicleRecordService,
	env *godotenv.Env, lg *slog.Logger, rts ...routes.Router) *Boot {
	return &Boot{event: event, eventConsumer: eventConsumer, retryConsumer: retryConsumer,
		rts: rts, rbt: rbt, cr: cr, retry: retry, vs: vs, lg: lg, env: env}
}

func (b *Boot) Boot(retry bool, retryFrom int64, retryLimit int) {
	env := godotenv.NewEnv()
	env.Load()
	if retry && retryLimit > 0 {
		manual := manual_retry.NewManualRetry(b.lg, b.env, b.vs)
		manual.FindResend(context.Background(), retryFrom, retryLimit)
	}
	r := gin.Default()
	for _, router := range b.rts {
		router.SetupRoutes(r)
	}
	docs.SwaggerInfo.BasePath = "/api"

	b.event.RegisterConsumer(b.eventConsumer)
	b.retry.RegisterConsumer(b.retryConsumer)
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
