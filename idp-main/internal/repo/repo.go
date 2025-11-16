package repo

import (
	healthzrepo "application/internal/repo/healthz"
	idprepo "application/internal/repo/idp"
	"application/internal/repo/vehicle"

	"github.com/google/wire"
)

var RepoProviderSet = wire.NewSet(idprepo.NewUserRepo, healthzrepo.NewHealthzDS, vehicle.NewVehicleRepo,
	NewDeviceRepository, NewWorkCalendarRepository,
	NewDistrictRepository,
)
