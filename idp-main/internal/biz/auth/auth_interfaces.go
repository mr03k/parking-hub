package authbiz

import "context"

type AutUsecaseInterface interface {
	DriverLogin(ctx context.Context, email, password string, driverID string) (string, string, error)
}
