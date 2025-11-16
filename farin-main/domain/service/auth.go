package service

import (
	"context"
	"errors"
	"farin/domain/entity"
	"farin/domain/repository"
	"farin/infrastructure/godotenv"
	"farin/util"
	"farin/util/encrypt"
	"github.com/golang-jwt/jwt"
	"github.com/mahdimehrabi/uploader/minio"
	"log/slog"
	"strconv"
	"time"
)

var ErrCalenderTimeInvalid = errors.New("calender time invalid")

type AuthService struct {
	env            *godotenv.Env
	logger         *slog.Logger
	userRepo       *repository.UserRepository
	driverRepo     *repository.DriverRepository
	contractorRepo *repository.ContractorRepository
	calenderRepo   *repository.CalenderRepository
	vehicleRepo    *repository.VehicleRepository
	deviceRepo     *repository.DeviceRepository
	assignmentRepo *repository.DriverAssignmentRepository
}

func NewAuthService(env *godotenv.Env, logger *slog.Logger, userRepository *repository.UserRepository,
	driverRepo *repository.DriverRepository, vehicleRepo *repository.VehicleRepository, deviceRepo *repository.DeviceRepository,
	assignmentRepo *repository.DriverAssignmentRepository, contractorRepo *repository.ContractorRepository) *AuthService {
	return &AuthService{
		env:            env,
		logger:         logger.With("layer", "AuthService"),
		userRepo:       userRepository,
		driverRepo:     driverRepo,
		vehicleRepo:    vehicleRepo,
		deviceRepo:     deviceRepo,
		assignmentRepo: assignmentRepo,
		contractorRepo: contractorRepo,
	}
}

func (s AuthService) CreateAccessToken(ctx context.Context, user *entity.User, exp int64, secret string) (string, error) {
	lg := s.logger.With("method", "CreateAccessToken")
	atClaims := jwt.MapClaims{
		"authorized":  true,
		"id":          user.ID,
		"email":       user.Email,
		"exp":         exp,
		"verified":    true,
		"roleID":      user.RoleID,
		"roleTitle":   user.Role.Title,
		"createdAt":   user.CreatedAt,
		"firstName":   user.FirstName,
		"lastName":    user.LastName,
		"phoneNumber": user.PhoneNumber,
	}
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(secret))
	if err != nil {
		lg.Error("Error in signing string", slog.Any("error", err))
		return "", err
	}
	return token, nil
}

func (s AuthService) CreateRefreshToken(ctx context.Context, user entity.User, exp int64, secret string) (string, error) {
	lg := s.logger.With("method", "CreateRefreshToken")
	atClaims := jwt.MapClaims{
		"authorized": true,
		"exp":        exp,
		"userID":     user.ID,
		"email":      user.Email,
	}
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(secret))
	if err != nil {
		lg.Error("Error in signing string", slog.Any("error", err))
		return "", err
	}
	return token, nil
}

func (s AuthService) CreateTokens(ctx context.Context, user *entity.User, remember bool) (map[string]string, error) {
	accessSecret := "access" + s.env.Secret
	expAccessToken := time.Now().Add(time.Minute * 15).Unix()
	accessToken, err := s.CreateAccessToken(ctx, user, expAccessToken, accessSecret)
	if err != nil {
		return nil, err
	}
	refreshSecret := "refresh" + s.env.Secret
	var expRefreshToken int64
	if remember {
		expRefreshToken = time.Now().Add(time.Hour * 24 * 15).Unix()
	} else {
		expRefreshToken = time.Now().Add(time.Hour * 24).Unix()
	}
	refreshToken, err := s.CreateRefreshToken(ctx, *user, expRefreshToken, refreshSecret)
	if err != nil {
		return nil, err
	}

	return map[string]string{
		"refreshToken":    refreshToken,
		"accessToken":     accessToken,
		"expRefreshToken": strconv.Itoa(int(expRefreshToken)),
		"expAccessToken":  strconv.Itoa(int(expAccessToken)),
	}, err
}

func (s AuthService) Login(ctx context.Context, email, enteredPassword string, remember bool) (*entity.User, map[string]string, error) {
	lg := s.logger.With("method", "Login")
	user, err := s.userRepo.GetByField(ctx, "phone_number", email)
	if errors.Is(err, repository.ErrUserNotFound) {
		return nil, nil, err
	}
	if err != nil {
		lg.Error("Error finding user", slog.Any("error", err))
		return nil, nil, err
	}
	user.ProfileImage = minio.GeneratePublicURL("http://", s.env.MinioHost,
		s.env.MinioProfilePictureBucket, user.ID)

	encryptedPassword := encrypt.HashSHA256(enteredPassword)
	if user.Password == encryptedPassword {
		tokensData, err := s.CreateTokens(ctx, user, remember)
		if err != nil {
			lg.Error("Failed to generate JWT tokens", slog.Any("error", err))
			return nil, nil, err
		}
		return user, tokensData, nil
	} else {
		return nil, nil, repository.ErrUserNotFound
	}
}

func (s AuthService) LoginDriver(ctx context.Context, user *entity.User, vehicleID string) (*entity.DriverAssignment,
	*entity.Contractor, []*entity.Device, error) {
	lg := s.logger.With("method", "Login")
	driver, err := s.driverRepo.GetByField(ctx, "user_id", user.ID)
	if err != nil {
		if errors.Is(err, repository.ErrDriverNotFound) {
			return nil, nil, nil, repository.ErrDriverNotFound
		}
		lg.Error("Error finding driver", slog.Any("error", err))
		return nil, nil, nil, err
	}
	var da *entity.DriverAssignment

	if driver != nil {
		da, err = s.assignmentRepo.GetByFields(ctx, driver.ID, vehicleID)
		if err != nil {
			if !errors.Is(err, repository.ErrDriverAssignmentNotFound) {
				lg.Error("Error finding driver", slog.Any("error", err))
				return nil, nil, nil, err
			}
		}
	}

	var contractor *entity.Contractor
	if driver != nil {
		contractor, err = s.contractorRepo.GetByField(ctx, "id", *driver.ContractorID)
	}

	var devices []*entity.Device
	if da != nil {
		devices, err = s.deviceRepo.GetMultipleByField(ctx, "vehicle_id", da.Vehicle.ID)
	}

	return da, contractor, devices, nil
}

func (s AuthService) RenewToken(ctx context.Context, refreshToken string) (accessToken string, expAccessToken int64, err error) {
	lg := s.logger.With("method", "RenewToken")
	var valid bool
	var atClaims jwt.MapClaims
	refreshSecret := "refresh" + s.env.Secret
	valid, atClaims, err = util.DecodeToken(refreshToken, refreshSecret)
	if err != nil {
		validationError := &jwt.ValidationError{}
		if errors.As(err, &validationError) {
			if validationError.Errors == jwt.ValidationErrorExpired {
				err = repository.ErrUserNotFound
				return "", 0, err
			}
		}
		lg.Error("Failed to decode refresh token", slog.Any("error", err))
		return "", 0, err
	}

	uid, ok := atClaims["userID"].(string)
	if !ok {
		err = repository.ErrUserNotFound
		lg.Error("Invalid token claims: missing userID", slog.Any("claims", atClaims))
		return "", 0, err
	}

	user, err := s.userRepo.GetByField(ctx, "id", uid)
	if errors.Is(err, repository.ErrUserNotFound) {
		return "", 0, err
	}
	if err != nil {
		lg.Error("Error finding user", slog.Any("error", err))
		return "", 0, err
	}

	if valid {
		expAccessToken := time.Now().Add(time.Minute * 30).Unix()
		accessToken, err := s.CreateAccessToken(ctx, user, expAccessToken, "access"+s.env.Secret)
		if err != nil {
			lg.Error("Error creating access token", slog.Any("error", err))
			return "", 0, err
		}
		return accessToken, expAccessToken, nil
	}
	err = repository.ErrUserNotFound
	return
}
