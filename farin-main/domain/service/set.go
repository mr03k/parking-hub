package service

import (
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	NewUserService,
	NewAuthService,
	NewContractorService,
	NewContractService,
	NewDriverService,
	NewVehicleService,
	NewDeviceService,
	NewVehicleRecordService, NewRoleService,
	NewCalenderService, NewRingService, NewDriverAssignmentService,
)
