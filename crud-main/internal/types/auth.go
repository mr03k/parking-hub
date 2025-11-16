package types

import (
	"git.abanppc.com/farin-project/crud/internal/model"
	"github.com/gin-gonic/gin"
)

// DriverLogin represents the request structure for login
type DriverLogin struct {
	VerifyRequest
}

// DriverLoginResponse represents the response structure for login
type DriverLoginResponse struct {
	User  model.Users `json:"user"`
	Role  model.Roles `json:"role"`
	Token string      `json:"token"`
}

// BindDriverLoginFromRequest binds and validates a DriverLogin request from Gin context
func BindDriverLoginFromRequest(c *gin.Context) (*DriverLogin, error) {
	var req DriverLogin
	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, err
	}
	return &req, nil
}
