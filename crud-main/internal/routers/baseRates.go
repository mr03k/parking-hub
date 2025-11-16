package routers

import (
	"github.com/gin-gonic/gin"

	"git.abanppc.com/farin-project/crud/internal/handler"
)

func init() {
	apiV1RouterFns = append(apiV1RouterFns, func(group *gin.RouterGroup) {
		baseRatesRouter(group, handler.NewBaseRatesHandler())
	})
}

func baseRatesRouter(group *gin.RouterGroup, h handler.BaseRatesHandler) {
	g := group.Group("/baseRates")

	// All the following routes use jwt authentication, you also can use middleware.Auth(middleware.WithVerify(fn))
	//g.Use(middleware.Auth())

	// If jwt authentication is not required for all routes, authentication middleware can be added
	// separately for only certain routes. In this case, g.Use(middleware.Auth()) above should not be used.

	g.POST("/", h.Create)          // [post] /api/v1/baseRates
	g.DELETE("/:id", h.DeleteByID) // [delete] /api/v1/baseRates/:id
	g.PUT("/:id", h.UpdateByID)    // [put] /api/v1/baseRates/:id
	g.GET("/:id", h.GetByID)       // [get] /api/v1/baseRates/:id
	g.POST("/list", h.List)        // [post] /api/v1/baseRates/list
}
