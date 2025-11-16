package authbiz

import (
	"context"
	"log/slog"

	"application/config"
	idpbiz "application/internal/biz/idp"

	"github.com/golang-jwt/jwt/v5"
)

type AuthUsecase struct {
	logger    *slog.Logger
	idpUC     idpbiz.UserUseCaseInterface
	jwtSecret []byte
}

type authConfig struct {
	JWT struct {
		Secret []byte
	}
}

func NewAuthUsecase(logger *slog.Logger, idpUC idpbiz.UserUseCaseInterface, cfg config.Config) AutUsecaseInterface {
	authConfig := new(authConfig)

	err := cfg.Unmarshal("service.auth", authConfig)
	if err != nil {
		logger.Error("unmarshal config", "error", err)
		panic(err)
	}

	logger.Debug("initAuthUsecase", "config", authConfig)
	return &AuthUsecase{
		logger:    logger.With("layer", "AuthUsecase"),
		idpUC:     idpUC,
		jwtSecret: authConfig.JWT.Secret,
	}
}

func (uc *AuthUsecase) DriverLogin(ctx context.Context, email, password, vehicleID string) (string, string, error) {
	logger := uc.logger.With("method", "DriverLogin")
	logger.Debug("DriverLogin", "email", email, "driverID", vehicleID)

	user, err := uc.idpUC.VerifyUser(ctx, email, password)
	if err != nil {
		logger.Error("error verify user", "error", err)
		return "", "", err
	}

	// for _, role := range roles {
	// 	if role.Name == "driver" {
	// 		return uc.jwtSecret, nil
	// 	}
	// }

	t := jwt.NewWithClaims(
		jwt.SigningMethodHS512,
		jwt.MapClaims{
			"iss":        "sharad.idp",
			"sub":        user.UserID,
			"msisdn":     user.Msisdn,
			"verified":   user.MsisdnVerified,
			"vehicle_id": vehicleID,
		})

	token, err := t.SignedString(uc.jwtSecret)
	if err != nil {
		logger.Error("error sign token", "failed", err)
		return "", "", err
	}

	return user.UserID, token, nil
}
