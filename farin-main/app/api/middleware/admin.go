package middleware

import (
	"farin/app/api/response"
	"farin/infrastructure/godotenv"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
)

// AuthMiddleware -> struct for authentication middleware
type AdminMiddleware struct {
	logger *slog.Logger
	env    *godotenv.Env
}

// NewAuthMiddleware -> new instance of AuthMiddleware
func NewAdminMiddleware(
	logger *slog.Logger,
	env *godotenv.Env,
) *AdminMiddleware {
	return &AdminMiddleware{
		logger: logger.With("Layer", "AuthMiddleware"),
		env:    env,
	}
}

func (m *AdminMiddleware) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := AuthenticatedUser(c)
		if user == nil {
			response.Custom(c, http.StatusUnauthorized, nil, "You must login to access this page 😥")
			c.Abort()
			return
		}

		if user.Role.Title != "Admin" {
			response.Custom(c, http.StatusForbidden, nil, "")
			c.Abort()
			return
		}

		c.Next()
	}
}
