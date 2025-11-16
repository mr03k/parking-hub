package handler

import (
	"context"
	"git.abanppc.com/farin-project/crud/internal/config"
	"git.abanppc.com/farin-project/crud/internal/ecode"
	"git.abanppc.com/farin-project/crud/internal/model"
	"git.abanppc.com/farin-project/crud/internal/types"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/zhufuyi/sponge/pkg/gin/middleware"
	"github.com/zhufuyi/sponge/pkg/gin/response"
	"github.com/zhufuyi/sponge/pkg/logger"
	"time"
)

var _ AuthHandler = (*authHandler)(nil)

// AuthHandler defining the handler interface
type AuthHandler interface {
	Login(c *gin.Context)
}

type authHandler struct {
}

// NewAuthHandler creates a new AuthHandler instance
func NewAuthHandler() AuthHandler {
	return &authHandler{}
}

// Login authenticates a user
// @Summary Login user
// @Description authenticate user and return a token
// @Tags authentication
// @accept json
// @Produce json
// @Param data body types.DriverLogin true "login information"
// @Param vehicle_id path string true "vehicle ID"
// @Success 200 {object} types.DriverLoginResponse{}
// @Router /api/auth/login/{vehicle_id} [post]
// @Security BearerAuth
func (h *authHandler) Login(c *gin.Context) {
	loginData := &types.DriverLogin{}
	err := c.ShouldBindJSON(loginData)
	if err != nil {
		logger.Warn("ShouldBindJSON error", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	vehicleID := c.Param("vehicle_id")

	ctx := middleware.WrapCtx(c)

	user, role, err := h.VerifyUser(ctx, loginData.Msisdn, loginData.Password)
	if err != nil {
		logger.Error("Verify user error", logger.Any("error", err))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	t := jwt.NewWithClaims(
		jwt.SigningMethodHS512,
		jwt.MapClaims{
			"iss":        "sharad.idp",
			"sub":        user.ID,
			"msisdn":     loginData.Msisdn,
			"verified":   true,
			"vehicle_id": vehicleID,
		})

	config := config.Get()

	token, err := t.SignedString([]byte(config.App.Secret))
	if err != nil {
		logger.Error("error sign token", logger.Any("error", err))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c, types.DriverLoginResponse{
		User:  *user,
		Role:  role,
		Token: token,
	})
}

func (h *authHandler) VerifyUser(ctx context.Context, msisdn, password string) (*model.Users, model.Roles, error) {

	return &model.Users{
			ID:           uuid.NewString(),
			Username:     "dsgdgsdsg",
			Address:      "SDG ds gsd gsd",
			Email:        "example@gmail.com",
			FirstName:    "dgsdgs",
			LastName:     "cvcvvc",
			NumberMobile: "+989120000000",
		}, model.Roles{
			uuid.NewString(), "driver", "sdggdsgds", time.Now().Unix()}, nil
}
