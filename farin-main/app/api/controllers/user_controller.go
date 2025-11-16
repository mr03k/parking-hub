package controller

import (
	"context"
	"farin/app/api/response"
	"farin/domain/dto"
	"farin/domain/service"
	"farin/infrastructure/godotenv"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log/slog"
	"strconv"
)

type UserController struct {
	userService *service.UserService
	logger      *slog.Logger
	env         *godotenv.Env
}

func NewUserController(logger *slog.Logger, us *service.UserService, env *godotenv.Env) *UserController {
	return &UserController{
		userService: us,
		logger:      logger.With("layer", "UserController"),
		env:         env,
	}
}

// CreateUser godoc
// @Summary      Create a new user
// @Description  Create a new user by providing user details
// @Tags         users
// @Accept       json
// @Security BearerAuth
// @Produce      json
// @Param        user  body      dto.UserRequest  true  "User Data"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success      201   {object}  response.Response[dto.UserResponse]
// @Failure      500   {object}  response.Response[swagger.EmptyObject]
// @Failure      400   {object}  response.Response[swagger.EmptyObject]
// @Router       /users [post]
func (uc *UserController) CreateUser(c *gin.Context) {
	lg := uc.logger.With("method", "CreateUser")
	var userRequest dto.UserRequest

	if err := c.ShouldBindJSON(&userRequest); err != nil {
		lg.Error("failed to bind user data", "error", err)
		response.BadRequest(c, "invalid request body")
		return
	}

	createdUser, err := uc.userService.CreateUser(context.Background(), userRequest.ToEntity(), userRequest)
	if err != nil {
		lg.Error("failed to create user", "error", err)
		response.InternalError(c)
		return
	}
	userResponse := dto.UserResponse{}
	userResponse.FromEntity(createdUser, uc.env)

	response.Created(c, userResponse)
}

// ListUsers godoc
// @Summary      List users
// @Description  Retrieve a list of users with optional filters and pagination
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        sortField  query     string  false  "Field to sort by"  default(created_at)
// @Param        sortOrder  query     string  false  "Sort order (asc/desc)"  default(asc)
// @Param        page       query     int     false  "Page number"  default(1)
// @Param        pageSize   query     int     false  "Page size"  default(10)
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success      200        {object}  response.Response[dto.UserListResponse]
// @Failure      500   {object}  response.Response[swagger.EmptyObject]
// @Failure      400   {object}  response.Response[swagger.EmptyObject]
// @Router       /users [get]
func (uc *UserController) ListUsers(c *gin.Context) {
	lg := uc.logger.With("method", "ListUsers")
	filters := make(map[string]interface{})
	sortField := c.DefaultQuery("sortField", "created_at")
	sortOrder := c.DefaultQuery("sortOrder", "asc")
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("pageSize", "10")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		response.BadRequest(c, "invalid page param")
	}
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		response.BadRequest(c, "invalid page size param")
	}

	users, total, err := uc.userService.ListUsers(context.Background(), filters, sortField, sortOrder, page, pageSize)
	if err != nil {
		lg.Error("failed to list users", "error", err)
		response.InternalError(c)
		return
	}
	userResponses := make([]dto.UserResponse, len(users))

	for i, user := range users {
		userResponse := dto.UserResponse{}
		userResponse.FromEntity(&user, uc.env)
		userResponses[i] = userResponse
	}

	response.Ok(c, dto.UserListResponse{
		Users: userResponses,
		Total: total,
	}, "")
}

// UpdateUser godoc
// @Summary      Update user details
// @Description  Update an existing user's details
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        user  body      dto.UserRequest  true  "Updated User Data"
// @Param        id  path      string  true  "id"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success      200   {object}  response.Response[dto.UserResponse]
// @Failure      500   {object}  response.Response[swagger.EmptyObject]
// @Failure      400   {object}  response.Response[swagger.EmptyObject]
// @Router       /users/{id} [put]
func (uc *UserController) UpdateUser(c *gin.Context) {
	lg := uc.logger.With("method", "UpdateUser")
	var userRequest dto.UserRequest

	if err := c.ShouldBindJSON(&userRequest); err != nil {
		lg.Error("failed to bind user data", "error", err)
		response.BadRequest(c, "invalid request body")
		return
	}

	user := userRequest.ToEntity()
	user.ID = c.Param("id")
	updatedUser, err := uc.userService.UpdateUser(context.Background(), user, userRequest)
	if err != nil {
		lg.Error("failed to update user", "error", err)
		response.InternalError(c)
		return
	}

	userResponse := dto.UserResponse{}
	userResponse.FromEntity(updatedUser, uc.env)

	response.Ok(c, userResponse, "")
}

// DeleteUser godoc
// @Summary      Delete a user
// @Description  Delete a user by ID
// @Tags         users
// @Param        id   path      string  true  "User ID"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success      200   {object}  response.Response[swagger.EmptyObject]
// @Failure      500   {object}  response.Response[swagger.EmptyObject]
// @Failure      400   {object}  response.Response[swagger.EmptyObject]
// @Failure      404   {object}  response.Response[swagger.EmptyObject]
// @Router       /users/{id} [delete]
func (uc *UserController) DeleteUser(c *gin.Context) {
	lg := uc.logger.With("method", "DeleteUser")
	id := c.Param("id")

	if _, err := uuid.Parse(id); err != nil {
		lg.Error("invalid user ID", "userID", id)
		response.BadRequest(c, "invalid request body")
		return
	}

	err := uc.userService.DeleteUser(context.Background(), id)
	if err != nil {
		lg.Error("failed to delete user", "error", err)
		response.InternalError(c)
		return
	}

	response.Ok(c, nil, "user deleted successfully")
}

// GetUserDetail godoc
// @Summary      Get user details
// @Description  Retrieve user details by ID
// @Tags         users
// @Param        id   path      string  true  "User ID"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success      200   {object}  response.Response[dto.UserResponse]
// @Failure      500   {object}  response.Response[swagger.EmptyObject]
// @Failure      404   {object}  response.Response[swagger.EmptyObject]
// @Router       /users/{id} [get]
func (uc *UserController) GetUserDetail(c *gin.Context) {
	lg := uc.logger.With("method", "GetUserDetail")
	id := c.Param("id")

	user, err := uc.userService.Detail(context.Background(), "id", id)
	if err != nil {
		lg.Error("failed to get user details", "error", err)
		response.NotFound(c)
		return
	}
	var userResponse dto.UserResponse
	userResponse.FromEntity(user, uc.env)

	response.Ok(c, userResponse, "")
}

// UploadUserPicture godoc
// @Summary      Upload user picture
// @Description  Upload a profile picture for a user
// @Tags         users
// @Accept       multipart/form-data
// @Produce      json
// @Param        id       path      string       true  "User ID"
// @Param        picture  formData  file         true  "Profile Picture"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success      200      {object}  response.Response[swagger.EmptyObject]
// @Failure      500   {object}  response.Response[swagger.EmptyObject]
// @Failure      400   {object}  response.Response[swagger.EmptyObject]
// @Failure      404   {object}  response.Response[swagger.EmptyObject]
// @Router       /users/{id}/picture [post]
func (uc *UserController) UploadUserPicture(c *gin.Context) {
	lg := uc.logger.With("method", "UploadUserPicture")
	id := c.Param("id")

	user, err := uc.userService.Detail(context.Background(), "id", id)
	if err != nil {
		lg.Error("user not found", "error", err)
		response.NotFound(c)
		return
	}

	file, _, err := c.Request.FormFile("picture")
	if err != nil {
		lg.Error("failed to retrieve picture", "error", err)
		response.BadRequest(c, "invalid request body")
		return
	}

	if err := uc.userService.UploadPicture(context.Background(), user, file); err != nil {
		lg.Error("failed to upload picture", "error", err)
		response.InternalError(c)
		return
	}

	response.Ok(c, nil, "picture uploaded successfully")
}
