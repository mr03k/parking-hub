package routes

import (
	"github.com/gin-gonic/gin"
)

type Router interface {
	SetupRoutes(engine *gin.Engine)
}

func CreateRouters(healthRouter *HealthRouter, vrRouter *VehicleRecordRouter) []Router {
	return []Router{
		healthRouter, vrRouter,
	}
}
