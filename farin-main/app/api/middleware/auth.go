package middleware

import (
	"farin/app/api/response"
	"farin/domain/entity"
	"farin/domain/repository"
	"farin/infrastructure/godotenv"
	"farin/util"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"log/slog"
	"net/http"
	"strings"
)

// AuthMiddleware -> struct for authentication middleware
type AuthMiddleware struct {
	logger   *slog.Logger
	env      *godotenv.Env
	userRepo *repository.UserRepository
}

// NewAuthMiddleware -> new instance of AuthMiddleware
func NewAuthMiddleware(
	logger *slog.Logger,
	env *godotenv.Env,
	userRepo *repository.UserRepository,
) *AuthMiddleware {
	return &AuthMiddleware{
		logger:   logger.With("Layer", "AuthMiddleware"),
		env:      env,
		userRepo: userRepo,
	}
}

type authHeader struct {
	Authorization string `header:"Authorization"`
}

func (m *AuthMiddleware) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		ah := authHeader{}
		if err := c.ShouldBindHeader(&ah); err == nil {
			strs := strings.Split(ah.Authorization, " ")
			if len(strs) != 2 {
				m.logger.Warn("Invalid access token format", slog.String("authorization", ah.Authorization))
				response.Custom(c, http.StatusUnauthorized, nil, "Your access token is not correct")
				c.Abort()
				return
			}

			bearer := strs[0]
			if bearer != "Bearer" {
				m.logger.Warn("Token does not start with 'Bearer'", slog.String("bearer", bearer))
				response.Custom(c, http.StatusUnauthorized, nil, "Your token doesn't start with 'Bearer '")
				c.Abort()
				return
			}

			accessToken := strs[1]
			valid, claims, err := util.DecodeToken(accessToken, "access"+m.env.Secret)
			if err != nil {
				m.logger.Error("Failed to decode token", slog.Any("error", err))
				response.Custom(c, http.StatusUnauthorized, nil, "You must login to access this page 😥")
				c.Abort()
				return
			}

			user, ok := m.claimsToUser(claims)
			if !ok {
				m.logger.Warn("Invalid user claims", slog.Any("claims", claims))
				response.Custom(c, http.StatusUnauthorized, nil, "You must login to access this page 😥")
				c.Abort()
				return
			}

			if valid {
				m.logger.Info("Token validated successfully", slog.String("userID", user.ID))
				c.Set("user", user)
				c.Next()
				return
			}
		}

		m.logger.Warn("Unauthorized access attempt")
		response.Custom(c, http.StatusUnauthorized, nil, "You must login to access this page 😥")
		c.Abort()
	}
}

// claimsToUser converts JWT claims to a user model
func (m AuthMiddleware) claimsToUser(claims jwt.MapClaims) (user *entity.User, ok bool) {
	defer func() {
		if r := recover(); r != nil {
			m.logger.Error("Recovered from panic while converting claims to user", slog.Any("panic", r))
			ok = false
		}
	}()

	user = new(entity.User)
	id, ok := claims["id"].(string)
	if !ok {
		m.logger.Warn("Missing userID in claims", slog.Any("claims", claims))
		return nil, false
	}

	user.ID = id
	user.Email = claims["email"].(string)
	user.FirstName = claims["firstName"].(string)
	user.LastName = claims["lastName"].(string)
	user.RoleID = claims["roleID"].(string)
	user.Role = entity.Role{
		Base: entity.Base{
			ID: claims["roleID"].(string),
		},
		Title: claims["roleTitle"].(string),
	}
	f := claims["createdAt"].(float64)
	user.CreatedAt = int64(f)
	ok = true
	return
}

// AuthenticatedUser returns the authenticated user from the Gin context (filled by JWT claims).
// If no user is stored in the Gin context, it returns nil.
func AuthenticatedUser(c *gin.Context) *entity.User {
	user := func() *entity.User {
		defer func() {
			if r := recover(); r != nil {
				return
			}
		}()
		return c.MustGet("user").(*entity.User)
	}()
	return user
}
