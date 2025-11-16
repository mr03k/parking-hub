package controller

import (
	"context"
	"errors"
	"farin/app/api/middleware"
	"farin/app/api/response"
	"farin/domain/dto"
	"farin/domain/repository"
	"farin/domain/service"
	"farin/infrastructure/godotenv"
	"farin/util"
	"farin/util/encrypt"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	logger      *slog.Logger
	env         *godotenv.Env
	userService *service.UserService
	authService *service.AuthService
}

func NewAuthController(
	env *godotenv.Env,
	authService *service.AuthService,
	userService *service.UserService,
	lg *slog.Logger,
) *AuthController {
	return &AuthController{
		env:         env,
		authService: authService,
		userService: userService,
		logger:      lg.With("layer", "AuthController"),
	}
}

// @Summary login
// @Schemes
// @Description jwt login
// @Tags auth
// @Accept json
// @Produce json
// @Param body body dto.LoginRequest true "request body"
// @Success 200 {object} response.Response[dto.LoginResponse]
// @failure 422 {object} response.Response[swagger.EmptyObject]
// @failure 401 {object} response.Response[swagger.EmptyObject]
// @Router /auth/login/ [post]
func (ac AuthController) Login(c *gin.Context) {
	// Create a context with a timeout of 5 seconds
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	// Data Parse
	var loginRequest dto.LoginRequest
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	user, tokensData, err := ac.authService.Login(ctx, loginRequest.PhoneNumber, loginRequest.Password, loginRequest.Remember)
	if errors.Is(err, repository.ErrUserNotFound) {
		response.BadRequest(c, "your login credentials is incorrect")
		return
	}

	if err != nil {
		response.InternalError(c)
		return
	}

	var loginResult dto.LoginResponse
	loginResult.AccessToken = tokensData["accessToken"]
	loginResult.RefreshToken = tokensData["refreshToken"]
	loginResult.ExpRefreshToken = tokensData["expRefreshToken"]
	loginResult.ExpAccessToken = tokensData["expAccessToken"]
	var userResponse dto.UserResponse
	userResponse.FromEntity(user, ac.env)
	loginResult.User = userResponse
	response.Ok(c, loginResult, "")
}

// @Summary driver-login
// @Schemes
// @Description jwt driver-login
// @Tags auth
// @Accept json
// @Produce json
// @Param body body dto.DriverLoginRequest true "request body"
// @Success 200 {object} response.Response[dto.DriverLoginResponse]
// @failure 422 {object} response.Response[swagger.EmptyObject]
// @failure 401 {object} response.Response[swagger.EmptyObject]
// @Router /auth/driver-login/ [post]
func (ac AuthController) DriverLogin(c *gin.Context) {
	// Create a context with a timeout of 5 seconds
	ctx, cancel := context.WithTimeout(c.Request.Context(), 55*time.Second)
	defer cancel()

	// Data Parse
	var loginRequest dto.DriverLoginRequest
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	user, tokensData, err := ac.authService.Login(ctx, loginRequest.PhoneNumber, loginRequest.Password, loginRequest.Remember)
	if errors.Is(err, repository.ErrUserNotFound) {
		response.BadRequest(c, "your login credentials is incorrect")
		return
	}
	if err != nil {
		response.InternalError(c)
		return
	}

	if user.Role.Title != "Driver" {
		response.BadRequest(c, "you are not driver")
		return
	}

	da, contractor, devices, err := ac.authService.LoginDriver(ctx, user, loginRequest.VehicleID)
	if err != nil {
		if errors.Is(err, repository.ErrDriverAssignmentNotFound) {
			response.BadRequest(c, "your login credentials is incorrect")
			return
		}
		response.InternalError(c)
		return
	}

	var loginResult dto.DriverLoginResponse
	var contractorResponse dto.ContractorResponse
	contractorResponse.FromEntity(contractor, ac.env)
	var calenderResponse dto.CalenderResponse
	calenderResponse.FromEntity(&da.Calender)
	var vehicleResponse dto.VehicleResponse
	vehicleResponse.FromEntity(&da.Vehicle)
	var driverResponse dto.DriverResponse
	driverResponse.FromEntity(&da.Driver)
	var ringResponse dto.RingResponse
	ringResponse.FromEntity(&da.Ring)
	devicesResponse := make([]dto.DeviceResponse, len(devices))
	for i, device := range devices {
		devicesResponse[i].FromEntity(ac.env, device)
	}

	loginResult.AccessToken = tokensData["accessToken"]
	loginResult.RefreshToken = tokensData["refreshToken"]
	loginResult.ExpRefreshToken = tokensData["expRefreshToken"]
	loginResult.ExpAccessToken = tokensData["expAccessToken"]
	loginResult.Contractor = contractorResponse
	loginResult.Calender = calenderResponse
	loginResult.Vehicle = vehicleResponse
	loginResult.Driver = driverResponse
	loginResult.Ring = ringResponse
	loginResult.Devices = devicesResponse
	var userResponse dto.UserResponse
	userResponse.FromEntity(user, ac.env)
	loginResult.User = userResponse
	response.Ok(c, loginResult, "")
}

// @Summary access token verify
// @Schemes
// @Description jwt access token verify
// @Tags auth
// @Accept json
// @Produce json
// @Param body body dto.AccessTokenReq true "request body"
// @Success 200 {object} response.Response[dto.LoginResponse]
// @failure 400 {object} response.Response[swagger.EmptyObject]
// @Router /auth/access-token-verify [post]
func (ac AuthController) AccessTokenVerify(c *gin.Context) {
	at := dto.AccessTokenReq{}
	if err := c.ShouldBindJSON(&at); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	accessToken := at.AccessToken
	accessSecret := "access" + ac.env.Secret
	valid, _, err := util.DecodeToken(accessToken, accessSecret)
	if err != nil {
		response.Custom(c, http.StatusBadRequest, gin.H{}, "access token is not valid")
		return
	}

	if valid {
		response.Ok(c, gin.H{}, "access token is valid")
	} else {
		response.Custom(c, http.StatusBadRequest, gin.H{}, "access token is not valid")
	}
}

// @Summary renew acces token
// @Schemes
// @Description jwt renew access token
// @Tags auth
// @Accept json
// @Produce json
// @Param body body dto.RefreshTokenRequest true "request body"
// @Success 200 {object} response.Response[dto.LoginResponse]
// @failure 400 {object} response.Response[swagger.EmptyObject]
// @Router /auth/renew-access-token [post]
func (ac AuthController) RenewToken(c *gin.Context) {
	// Create a context with a timeout of 5 seconds
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	rtr := dto.RefreshTokenRequest{}
	if err := c.ShouldBindJSON(&rtr); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	accessToken, exp, err := ac.authService.RenewToken(ctx, rtr.RefreshToken)
	if errors.Is(err, repository.ErrUserNotFound) {
		response.Custom(c, http.StatusBadRequest, gin.H{}, "access token is not valid")
		return
	}
	if err != nil {
		response.InternalError(c)
		return
	}
	response.Ok(c, dto.AccessTokenRes{AccessToken: accessToken, ExpAccessToken: strconv.Itoa(int(exp))}, "")
}

// @Summary set password
// @Schemes
// @Description set password
// @Tags user,auth
// @Accept json
// @Param body body dto.PasswordRequest true "request body"
// @Param Authorization header string true "bearer authentication"
// @Produce json
// @param Authorization header string true "Authorization"
// @Param        contract  body      dto.ContractRequest  true  "Contract Data"
// @failure 200 {object} response.Response[swagger.EmptyObject]
// @failure 400 {object} response.Response[swagger.EmptyObject]
// @failure 401 {object} response.Response[swagger.EmptyObject]
// @Router /auth/password/ [PUT]
func (uc *AuthController) UpdatePassword(c *gin.Context) {
	// Create a context with timeout
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	var req dto.PasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	user := middleware.AuthenticatedUser(c)
	if user == nil {
		response.Custom(c, http.StatusUnauthorized, nil, "You must login to access this page 😥")
		return
	}
	user.Password = encrypt.HashSHA256(req.Password)

	if _, err := uc.userService.UpdateUser(ctx, user, dto.UserRequest{}); err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			response.NotFound(c)
		} else if errors.Is(err, repository.ErrUserNotFound) {
			response.BadRequest(c, err.Error())
			return
		} else {
			response.InternalError(c)
		}
		return
	}

	response.Ok(c, nil, "your password has been updated")
}
