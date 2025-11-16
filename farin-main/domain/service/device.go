package service

import (
	"bytes"
	"context"
	"errors"
	"farin/domain/dto"
	"farin/domain/entity"
	"farin/domain/repository"
	"farin/infrastructure/godotenv"
	"fmt"
	"github.com/google/uuid"
	"github.com/mahdimehrabi/uploader"
	"log/slog"
	"strings"
)

type DeviceService struct {
	logger     *slog.Logger
	deviceRepo *repository.DeviceRepository
	fr         uploader.FileRepository
	env        *godotenv.Env
}

func NewDeviceService(logger *slog.Logger, deviceRepo *repository.DeviceRepository, fr uploader.FileRepository,
	env *godotenv.Env) *DeviceService {
	return &DeviceService{
		logger:     logger.With("layer", "DeviceService"),
		deviceRepo: deviceRepo,
		fr:         fr,
		env:        env,
	}
}
func (s *DeviceService) CreateDevice(ctx context.Context, device *entity.Device, deviceDTO *dto.DeviceRequest) (*entity.Device, error) {
	lg := s.logger.With("method", "CreateDevice")
	device.ID = uuid.NewString()

	var errCallback func()

	if deviceDTO.ImageInsurance != nil {
		rs := bytes.NewReader(deviceDTO.ImageInsurance)
		err := s.fr.UploadPublicFile(ctx, s.env.MinioProfilePictureBucket, device.ID+"_image_insurance", "application/octet-stream", rs)
		if err != nil {
			lg.Error("failed to upload picture", slog.Any("error", err))
			return nil, err
		}
		device.ImageInsurance = fmt.Sprintf("%s/%s", s.env.MinioProfilePictureBucket, device.ID+"_image_insurance")

		errCallback = func() {
			if err1 := s.deleteFile(ctx, device.ImageInsurance, lg); err1 != nil {
				slog.Error("failed to delete file", slog.Any("error", err1))
			}
		}
	}

	if deviceDTO.ImageContract != nil {
		rs := bytes.NewReader(deviceDTO.ImageContract)
		err := s.fr.UploadPublicFile(ctx, s.env.MinioProfilePictureBucket, device.ID+"_image_contract", "application/octet-stream", rs)
		if err != nil {
			lg.Error("failed to upload picture", slog.Any("error", err))
			if errCallback != nil {
				errCallback()
			}
			return nil, err
		}
		device.ImageContract = fmt.Sprintf("%s/%s", s.env.MinioProfilePictureBucket, device.ID+"_image_contract")

		// Store the previous callback
		previousCallback := errCallback
		errCallback = func() {
			if previousCallback != nil {
				previousCallback()
			}
			if err1 := s.deleteFile(ctx, device.ImageContract, lg); err1 != nil {
				slog.Error("failed to delete file", slog.Any("error", err1))
			}
		}
	}

	createdDevice, err := s.deviceRepo.Create(ctx, device)
	if err != nil {
		lg.Error("failed to create device", "error", err.Error())
		if errCallback != nil {
			errCallback()
		}
		return nil, err
	}
	lg.Info("device created", "deviceID", device.ID)
	return createdDevice, nil
}

func (s DeviceService) deleteFile(ctx context.Context, file string, lg *slog.Logger) error {
	if file == "" {
		return nil
	}
	strs := strings.Split(file, "/")
	if len(strs) < 2 {
		lg.Warn("failed to delete file", slog.String("file", file))
		return errors.New("failed to delete file the file name is out of format")
	}
	if err := s.fr.DeleteFile(ctx, strs[0], strs[1]); err != nil {
		lg.Error("failed to delete file", slog.String("file", file), slog.Any("error", err))
		return err
	}
	return nil
}

func (s *DeviceService) ListDevices(ctx context.Context, filters map[string]interface{}, sortField, sortOrder string, page, pageSize int) ([]entity.Device, int64, error) {
	logger := s.logger.With("method", "ListDevices")
	devices, total, err := s.deviceRepo.List(ctx, filters, sortField, sortOrder, page, pageSize)
	if err != nil {
		logger.Error("failed to list devices", "error", err.Error())
		return nil, 0, err
	}
	logger.Info("devices listed", "totalDevices", total)
	return devices, total, nil
}

func (s *DeviceService) UpdateDevice(ctx context.Context, device *entity.Device, deviceDTO dto.DeviceRequest) (*entity.Device, error) {
	lg := s.logger.With("method", "UpdateDevice")
	existingDevice, err := s.deviceRepo.GetByField(ctx, "id", device.ID)
	if err != nil {
		if errors.Is(err, repository.ErrDeviceNotFound) {
			lg.Warn("device not found for update", "deviceID", device.ID)
			return nil, repository.ErrDeviceNotFound
		}
		lg.Error("failed to get device for update", "error", err.Error())
		return nil, err
	}

	var errCallback func()
	if deviceDTO.ImageInsurance != nil {
		rs := bytes.NewReader(deviceDTO.ImageInsurance)
		err := s.fr.UploadPublicFile(ctx, s.env.MinioProfilePictureBucket, device.ID+"_image_insurance", "application/octet-stream", rs)
		if err != nil {
			lg.Error("failed to upload picture", slog.Any("error", err))
			return nil, err
		}
		device.ImageInsurance = fmt.Sprintf("%s/%s", s.env.MinioProfilePictureBucket, device.ID+"_image_insurance")

		errCallback = func() {
			if err1 := s.deleteFile(ctx, device.ImageInsurance, lg); err1 != nil {
				slog.Error("failed to delete file", slog.Any("error", err1))
			}
		}
	}

	if deviceDTO.ImageContract != nil {
		rs := bytes.NewReader(deviceDTO.ImageContract)
		err := s.fr.UploadPublicFile(ctx, s.env.MinioProfilePictureBucket, device.ID+"_image_contract", "application/octet-stream", rs)
		if err != nil {
			lg.Error("failed to upload picture", slog.Any("error", err))
			return nil, err
		}
		device.ImageContract = fmt.Sprintf("%s/%s", s.env.MinioProfilePictureBucket, device.ID+"_image_contract")

		errCallback = func() {
			errCallback()
			if err1 := s.deleteFile(ctx, device.ImageContract, lg); err1 != nil {
				slog.Error("failed to delete file", slog.Any("error", err1))
			}
		}
	}

	lg.Info("updating device", "deviceID", existingDevice.ID)
	_, err = s.deviceRepo.Update(ctx, device)
	if err != nil {
		lg.Error("failed to update device", "error", err.Error())
		errCallback()
		return nil, err
	}
	updatedDevice, err := s.deviceRepo.GetByField(ctx, "id", existingDevice.ID)
	if err != nil {
		lg.Error("failed to  findupdate device", "error", err.Error())
		return nil, err
	}
	lg.Info("device updated", "deviceID", updatedDevice.ID)
	return updatedDevice, nil
}

func (s *DeviceService) DeleteDevice(ctx context.Context, id string) error {
	logger := s.logger.With("method", "DeleteDevice")
	existingDevice, err := s.deviceRepo.GetByField(ctx, "id", id)
	if err != nil {
		if errors.Is(err, repository.ErrDeviceNotFound) {
			logger.Warn("device not found for deletion", "deviceID", id)
			return nil
		}
		logger.Error("failed to find device for deletion", "error", err.Error())
		return err
	}

	logger.Info("deleting device", "deviceID", existingDevice.ID)
	err = s.deviceRepo.Delete(ctx, id)
	if err != nil {
		logger.Error("failed to delete device", "error", err.Error())
		return err
	}
	logger.Info("device deleted", "deviceID", existingDevice.ID)
	return nil
}

func (s *DeviceService) Detail(ctx context.Context, id, value string) (*entity.Device, error) {
	logger := s.logger.With("method", "Detail")
	device, err := s.deviceRepo.GetByField(ctx, id, value)
	if err != nil {
		if errors.Is(err, repository.ErrDeviceNotFound) {
			logger.Warn("device not found for detail", "field", id, "value", value)
			return nil, errors.New("device not found")
		}
		logger.Error("failed to get device details", "error", err.Error())
		return nil, err
	}
	logger.Info("device details retrieved", "deviceID", device.ID)
	return device, nil
}
