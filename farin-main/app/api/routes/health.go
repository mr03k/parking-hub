package routes

import (
	controller "farin/app/api/controllers"
	"github.com/gin-gonic/gin"
)

type HealthRouter struct {
	healthController *controller.HealthController
}

func NewHealthRouter(healthController *controller.HealthController) *HealthRouter {
	return &HealthRouter{healthController: healthController}
}

func (rh *HealthRouter) SetupRoutes(router *gin.Engine) {
	g := router.Group("/")
	{
		g.GET("/liveness", rh.healthController.Liveness)
		g.GET("/readiness", rh.healthController.Readiness)
	}
}
