package healthz

import (
	"context"
	"git.abanppc.com/farin-project/event/pkg/mongo"
	"git.abanppc.com/farin-project/event/pkg/rabbit"
	"go.opentelemetry.io/otel"
	"log"
	"log/slog"
	"os"
	"time"
)

type Healthz struct {
	logger    *slog.Logger
	lastReady time.Time
	rbt       *rabbit.Rabbit
	mongo     *mongo.Mongo
}

func NewHealthz(logger *slog.Logger, rbt *rabbit.Rabbit, mongo *mongo.Mongo) *Healthz {
	return &Healthz{
		logger:    logger.With("layer", "MuxHealthzService"),
		lastReady: time.Now(),
		rbt:       rbt,
		mongo:     mongo,
	}
}

// Healthz Readiness
func (s *Healthz) HealthzReadiness() {
	ctx, span := otel.Tracer("healthz").Start(context.Background(), "readiness")
	defer span.End()

	logger := s.logger.With("method", "HealthzReadiness", "ctx", ctx)
	if s.lastReady.Compare(time.Now().Add(-(time.Minute * 5))) == -1 {
		logger.Error("healthz readiness timeout,killing the instance...")
		time.Sleep(3 * time.Second)
		log.Fatal("healthz readiness timeout,killing the instance...")
	}

	if err := s.rbt.HealthCheck(); err != nil {
		logger.Warn("error checking rabbitMQ health", "error", err)
		if err := s.rbt.Setup(); err != nil {
			logger.Warn("error setting up rabbitMQ", "error", err)
		}
		return
	}
	if err := s.mongo.HealthCheck(); err != nil {
		logger.Warn("error checking mongodb health", "error", err)
		return
	}

	f, err := os.OpenFile("readiness", os.O_RDONLY|os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		logger.Error("healthz readiness create file failed,killing the instance...")
		time.Sleep(3 * time.Second)
		log.Fatal("healthz readiness create file failed,killing the instance...")
	}
	_, err = f.Write([]byte("ok"))
	if err != nil {
		logger.Error("healthz readiness create file failed,killing the instance...")
		time.Sleep(3 * time.Second)
		log.Fatal("healthz readiness create file failed,killing the instance...")
	}
	s.lastReady = time.Now()
}
