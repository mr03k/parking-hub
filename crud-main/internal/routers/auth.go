package routers

import (
	"github.com/gin-gonic/gin"

	"git.abanppc.com/farin-project/crud/internal/handler"
)

func init() {
	apiV1RouterFns = append(apiV1RouterFns, func(group *gin.RouterGroup) {
		authsRouter(group, handler.NewAuthHandler())
	})
}

func authsRouter(group *gin.RouterGroup, h handler.AuthHandler) {
	g := group.Group("/auth")

	// All the following routes use jwt authentication, you also can use middleware.Auth(middleware.WithVerify(fn))
	//g.Use(middleware.Auth())

	// If jwt authentication is not required for all routes, authentication middleware can be added
	// separately for only certain routes. In this case, g.Use(middleware.Auth()) above should not be used.

	g.POST("/login/:vehicle_id", h.Login)
}
