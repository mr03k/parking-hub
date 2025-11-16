package repository

import (
	"github.com/google/wire"
	"github.com/mahdimehrabi/uploader"
	"github.com/mahdimehrabi/uploader/minio"
)

var ProviderSet = wire.NewSet(
	wire.Bind(new(uploader.FileRepository), new(*minio.MinIOFileRepository)),
	minio.NewMinIOFileRepository,
	NewUserRepository,
	NewContractorRepository,
	NewContractRepository,
	NewDriverRepository,
	NewVehicleRepository, NewRoleRepository,
	NewDeviceRepository, NewEventRecordRabbitMQ,
	NewCalenderRepository, NewVehicleRecordRepository, NewCitizenVehiclePhotoRepository,
	NewRingRepository, NewDriverAssignmentRepository,
)
