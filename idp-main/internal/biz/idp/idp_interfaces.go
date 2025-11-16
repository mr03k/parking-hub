package idpbiz

import (
	"context"

	idpentities "application/internal/entity/idp"
)

type UserUseCaseInterface interface {
	// create user and return user id and user and error
	CreateUser(ctx context.Context, user *idpentities.User) (string, *idpentities.User, error)
	// List of users
	ListUser(ctx context.Context) ([]idpentities.User, error)

	// Get User Only By ID
	// id can be uuid
	GetUserByID(ctx context.Context, id string) (*idpentities.User, error)

	// Delete user
	DeleteUser(ctx context.Context, id string) error

	// Verify User
	VerifyUser(ctx context.Context, msisdn, password string) (*idpentities.User, error)

	UpdateUser(ctx context.Context, userID string, updatedUser *idpentities.User) error
}

type UserRepositoryInterface interface {
	// Create User and return userID and user and error
	CreateUser(ctx context.Context, user *idpentities.User) (string, *idpentities.User, error)
	// List of User
	ListUser(ctx context.Context) ([]idpentities.User, error)

	// Get User Only By ID
	// id can be email or msisdn or uuid
	GetUserByID(ctx context.Context, id string) (*idpentities.User, error)

	// Delete user
	DeleteUser(ctx context.Context, id string) error

	GetUserByMsisdn(ctx context.Context, msisdn string) (*idpentities.User, error)

	UpdateUser(ctx context.Context, userID string, updatedUser *idpentities.User) error
}
