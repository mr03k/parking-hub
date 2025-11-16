package controller

import "github.com/google/wire"

var ProviderSet = wire.NewSet(NewAuthController, NewHealthController,
	NewUserController, NewContractController, NewVehicleController, NewDeviceController,
	NewContractorController, NewDriverController, NewCalenderController, NewRingController,
	NewDriverAssignmentController, NewRoleController)
