package controller

import (
	"context"
	"git.abanppc.com/farin-project/state/app/api/response"
	"git.abanppc.com/farin-project/state/app/rabbit/consumers"
	gormdb "git.abanppc.com/farin-project/state/infrastructure/gorm"
	"git.abanppc.com/farin-project/state/infrastructure/rabbit"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"time"
)

type HealthController struct {
	lastReady     time.Time
	rbt           *rabbit.Rabbit
	logger        *slog.Logger
	gorm          *gormdb.GORMDB
	eventConsumer *consumers.StateConsumer
}

func NewHealthController(logger *slog.Logger, rbt *rabbit.Rabbit,
	gorm *gormdb.GORMDB, eventConsumer *consumers.StateConsumer) *HealthController {
	return &HealthController{rbt: rbt, gorm: gorm,
		lastReady:     time.Now(),
		logger:        logger.With("layer", "HealthController"),
		eventConsumer: eventConsumer,
	}
}

func (hc *HealthController) Liveness(c *gin.Context) {
	lg := hc.logger.With("method", "Liveness")

	if hc.lastReady.Compare(time.Now().Add(-(time.Minute * 5))) == -1 {
		lg.Warn("5 minutes of Radiness failure,making liveness fail to restart...")
		response.Custom(c, http.StatusServiceUnavailable, nil, "")
		return
	}

	response.Ok(c, nil, "")
}

func (hc *HealthController) Readiness(c *gin.Context) {
	lg := hc.logger.With("method", "Readiness")

	if err := hc.rbt.HealthCheck(); err != nil {
		lg.Error("rabbitMQ health check failed", "error", err)
		response.Custom(c, http.StatusServiceUnavailable, nil, "")
	}

	if err := hc.gorm.HealthCheck(context.Background()); err != nil {
		lg.Error("mongodb health check failed", "error", err)
		response.Custom(c, http.StatusServiceUnavailable, nil, "")
	}

	hc.lastReady = time.Now()

	response.Ok(c, nil, "")
}
