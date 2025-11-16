package routes

import (
	controller "farin/app/api/controllers"
	"farin/app/api/middleware"
	"github.com/gin-gonic/gin"
)

type AuthRouter struct {
	authController  *controller.AuthController
	authMiddleware  *middleware.AuthMiddleware
	adminMiddleware *middleware.AdminMiddleware
}

func NewAuthRouter(authController *controller.AuthController,
	authMiddleware *middleware.AuthMiddleware, adminMiddleware *middleware.AdminMiddleware) *AuthRouter {

	return &AuthRouter{authController: authController,
		authMiddleware:  authMiddleware,
		adminMiddleware: adminMiddleware}
}

func (rh *AuthRouter) SetupRoutes(router *gin.Engine) {
	g := router.Group("/api/auth")
	{
		g.POST("/login", rh.authController.Login)
		g.POST("/driver-login", rh.authController.DriverLogin)
		g.POST("/access-token-verify", rh.authController.AccessTokenVerify)
		g.POST("/renew-access-token", rh.authController.RenewToken)
		g.Use(rh.authMiddleware.Handle()).Use(rh.adminMiddleware.Handle()).PUT("/password", rh.authController.UpdatePassword)
	}
}
