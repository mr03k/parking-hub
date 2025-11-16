package controller

import (
	"context"
	"farin/app/api/response"
	"farin/app/rabbit/consumers"
	"farin/domain/service"
	gormdb "farin/infrastructure/gorm"
	"farin/infrastructure/rabbit"
	"github.com/gin-gonic/gin"
	"github.com/mahdimehrabi/uploader/minio"
	"log/slog"
	"net/http"
	"time"
)

type HealthController struct {
	service       service.UserService
	lastReady     time.Time
	rbt           *rabbit.Rabbit
	logger        *slog.Logger
	minio         *minio.Minio
	gorm          *gormdb.GORMDB
	eventConsumer *consumers.EventConsumer
}

func NewHealthController(logger *slog.Logger, rbt *rabbit.Rabbit, minio *minio.Minio,
	gorm *gormdb.GORMDB, eventConsumer *consumers.EventConsumer) *HealthController {
	return &HealthController{rbt: rbt, gorm: gorm,
		lastReady:     time.Now(),
		minio:         minio,
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

	if err := hc.minio.ReadinessCheck(); err != nil {
		lg.Error("minio health check failed", "error", err)
		response.Custom(c, http.StatusServiceUnavailable, nil, "")
	}
	hc.lastReady = time.Now()

	response.Ok(c, nil, "")
}
