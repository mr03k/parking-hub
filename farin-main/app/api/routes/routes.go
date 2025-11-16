package routes

import (
	"github.com/gin-gonic/gin"
)

type Router interface {
	SetupRoutes(engine *gin.Engine)
}

func CreateRouters(authRouter *AuthRouter, healthRouter *HealthRouter, userRouter *UserRouter) []Router {
	return []Router{
		authRouter, healthRouter, userRouter,
	}
}
