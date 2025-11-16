package idpbiz

import (
	idpentities "application/internal/entity/idp"
	"context"
	"errors"
	"log/slog"

	"golang.org/x/crypto/bcrypt"
)

type UserUsecase struct {
	logger *slog.Logger
	repo   UserRepositoryInterface
}

func NewUserUseCase(
	logger *slog.Logger,
	repo UserRepositoryInterface,
) UserUseCaseInterface {
	return &UserUsecase{
		logger: logger.With("layer", "IDPUseCase"),
		repo:   repo,
	}
}

func (u *UserUsecase) CreateUser(ctx context.Context, user *idpentities.User) (string, *idpentities.User, error) {
	logger := u.logger.With("method", "CreateUser")
	logger.Debug("usecase CreateUser")

	if user.ClearPassword != nil {
		var err error
		user.EncryptedPassword, err = bcrypt.GenerateFromPassword([]byte(user.ClearPassword), bcrypt.DefaultCost)
		if err != nil {
			logger.Error("error bcrypt.GenerateFromPassword", "error", err)
			return "", nil, err
		}
	}

	logger.Debug("usecase CreateUser", "user", user)

	id, encuser, err := u.repo.CreateUser(ctx, user)
	if err != nil {
		logger.Error("error CreateUser", "error", err)
		return "", nil, err
	}

	return id, encuser, nil
}

func (u *UserUsecase) ListUser(ctx context.Context) ([]idpentities.User, error) {
	logger := u.logger.With("method", "ListUsers")
	logger.Debug("usecase ListUsers")

	users, err := u.repo.ListUser(ctx)
	if err != nil {
		logger.Error("error ListUsers", "error", err)
		return nil, err
	}

	return users, nil
}

func (u *UserUsecase) GetUserByID(ctx context.Context, id string) (*idpentities.User, error) {
	logger := u.logger.With("method", "GetUserByID")
	logger.Debug("usecase GetUserByID", "id", id)

	user, err := u.repo.GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrorUserNotFount) {
			return nil, err
		}
		logger.Error("error GetUserByID", "error", err)
		return nil, err
	}

	return user, nil
}

func (u *UserUsecase) DeleteUser(ctx context.Context, id string) error {
	logger := u.logger.With("method", "DeleteUser")
	logger.Debug("usecase DeleteUser", "id", id)

	err := u.repo.DeleteUser(ctx, id)
	if err != nil {
		logger.Error("error DeleteUser", "error", err)
		return err
	}
	return nil
}

func (u *UserUsecase) UpdateUser(ctx context.Context, userID string, updatedUser *idpentities.User) error {
	logger := u.logger.With("method", "UpdateUser")
	logger.Debug("usecase UpdateUser", "userID", userID, "updatedUser", updatedUser)

	// Retrieve existing user to ensure it exists and handle not found errors
	existingUser, err := u.repo.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, ErrorUserNotFount) {
			logger.Error("user not found", "userID", userID, "error", err)
			return ErrorUserNotFount
		}
		logger.Error("error retrieving user for update", "userID", userID, "error", err)
		return err
	}

	// Handle password encryption if a new password is provided
	if updatedUser.ClearPassword != nil {
		updatedUser.EncryptedPassword, err = bcrypt.GenerateFromPassword(updatedUser.ClearPassword, bcrypt.DefaultCost)
		if err != nil {
			logger.Error("error encrypting password", "error", err)
			return err
		}
	}

	// Merge new data into the existing user
	if updatedUser.Msisdn != nil {
		existingUser.Msisdn = updatedUser.Msisdn
	}
	existingUser.MsisdnVerified = updatedUser.MsisdnVerified
	if updatedUser.EncryptedPassword != nil {
		existingUser.EncryptedPassword = updatedUser.EncryptedPassword
	}

	// Update the user in the repository
	err = u.repo.UpdateUser(ctx, userID, existingUser)
	if err != nil {
		logger.Error("error updating user in repository", "userID", userID, "error", err)
		return err
	}

	logger.Debug("usecase UpdateUser successful", "userID", userID)
	return nil
}

func (u *UserUsecase) VerifyUser(ctx context.Context, msisdn, password string) (*idpentities.User, error) {
	logger := u.logger.With("method", "VerifyUser")
	logger.Debug("usecase VerifyUser", "msisdn", msisdn)

	user, err := u.repo.GetUserByMsisdn(ctx, msisdn)
	if err != nil {
		logger.Error("error VerifyUser", "error", err)
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.EncryptedPassword), []byte(password))
	if err != nil {
		logger.Error("error VerifyUser", "error", err)
		return nil, ErrorValidationFailed
	}

	return user, nil
}
