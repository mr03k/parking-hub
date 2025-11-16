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

type VehicleService struct {
	logger      *slog.Logger
	vehicleRepo *repository.VehicleRepository
	fr          uploader.FileRepository
	env         *godotenv.Env
}

func NewVehicleService(logger *slog.Logger, vehicleRepo *repository.VehicleRepository, fr uploader.FileRepository,
	env *godotenv.Env) *VehicleService {
	return &VehicleService{
		logger:      logger.With("layer", "VehicleService"),
		vehicleRepo: vehicleRepo,
		fr:          fr,
		env:         env,
	}
}
func (s *VehicleService) CreateVehicle(ctx context.Context, vehicle *entity.Vehicle, vehicleDTO *dto.VehicleRequest) (*entity.Vehicle, error) {
	lg := s.logger.With("method", "CreateVehicle")
	vehicle.ID = uuid.NewString()

	var errCallback func()
	if vehicleDTO.ImageDocumentVehicle != nil {
		rs := bytes.NewReader(vehicleDTO.ImageDocumentVehicle)
		err := s.fr.UploadPublicFile(ctx, s.env.MinioProfilePictureBucket, vehicle.ID+"_image_document_vehicle", "application/octet-stream", rs)
		if err != nil {
			lg.Error("failed to upload document image", slog.Any("error", err))
			return nil, err
		}
		vehicle.ImageDocumentVehicle = fmt.Sprintf("%s/%s", s.env.MinioProfilePictureBucket, vehicle.ID+"_image_document_vehicle")

		errCallback = func() {
			if err1 := s.deleteFile(ctx, vehicle.ImageDocumentVehicle, lg); err1 != nil {
				slog.Error("failed to delete file", slog.Any("error", err1))
			}
		}
	}

	if vehicleDTO.ImageCardVehicle != nil {
		rs := bytes.NewReader(vehicleDTO.ImageCardVehicle)
		err := s.fr.UploadPublicFile(ctx, s.env.MinioProfilePictureBucket, vehicle.ID+"_image_card_vehicle", "application/octet-stream", rs)
		if err != nil {
			lg.Error("failed to upload card image", slog.Any("error", err))
			if errCallback != nil {
				errCallback()
			}
			return nil, err
		}
		vehicle.ImageCardVehicle = fmt.Sprintf("%s/%s", s.env.MinioProfilePictureBucket, vehicle.ID+"_image_card_vehicle")

		if errCallback != nil {
			oldCallback := errCallback
			errCallback = func() {
				oldCallback()
				if err1 := s.deleteFile(ctx, vehicle.ImageCardVehicle, lg); err1 != nil {
					slog.Error("failed to delete file", slog.Any("error", err1))
				}
			}
		} else {
			errCallback = func() {
				if err1 := s.deleteFile(ctx, vehicle.ImageCardVehicle, lg); err1 != nil {
					slog.Error("failed to delete file", slog.Any("error", err1))
				}
			}
		}
	}

	if vehicleDTO.ThirdPartyInsuranceImage != nil {
		rs := bytes.NewReader(vehicleDTO.ThirdPartyInsuranceImage)
		err := s.fr.UploadPublicFile(ctx, s.env.MinioProfilePictureBucket, vehicle.ID+"_third_party_insurance", "application/octet-stream", rs)
		if err != nil {
			lg.Error("failed to upload third party insurance image", slog.Any("error", err))
			if errCallback != nil {
				errCallback()
			}
			return nil, err
		}
		vehicle.ThirdPartyInsuranceImage = fmt.Sprintf("%s/%s", s.env.MinioProfilePictureBucket, vehicle.ID+"_third_party_insurance")

		if errCallback != nil {
			oldCallback := errCallback
			errCallback = func() {
				oldCallback()
				if err1 := s.deleteFile(ctx, vehicle.ThirdPartyInsuranceImage, lg); err1 != nil {
					slog.Error("failed to delete file", slog.Any("error", err1))
				}
			}
		} else {
			errCallback = func() {
				if err1 := s.deleteFile(ctx, vehicle.ThirdPartyInsuranceImage, lg); err1 != nil {
					slog.Error("failed to delete file", slog.Any("error", err1))
				}
			}
		}
	}

	if vehicleDTO.BodyInsuranceImage != nil {
		rs := bytes.NewReader(vehicleDTO.BodyInsuranceImage)
		err := s.fr.UploadPublicFile(ctx, s.env.MinioProfilePictureBucket, vehicle.ID+"_body_insurance", "application/octet-stream", rs)
		if err != nil {
			lg.Error("failed to upload body insurance image", slog.Any("error", err))
			if errCallback != nil {
				errCallback()
			}
			return nil, err
		}
		vehicle.BodyInsuranceImage = fmt.Sprintf("%s/%s", s.env.MinioProfilePictureBucket, vehicle.ID+"_body_insurance")

		if errCallback != nil {
			oldCallback := errCallback
			errCallback = func() {
				oldCallback()
				if err1 := s.deleteFile(ctx, vehicle.BodyInsuranceImage, lg); err1 != nil {
					slog.Error("failed to delete file", slog.Any("error", err1))
				}
			}
		} else {
			errCallback = func() {
				if err1 := s.deleteFile(ctx, vehicle.BodyInsuranceImage, lg); err1 != nil {
					slog.Error("failed to delete file", slog.Any("error", err1))
				}
			}
		}
	}

	createdVehicle, err := s.vehicleRepo.Create(ctx, vehicle)
	if err != nil {
		lg.Error("failed to create vehicle", "error", err.Error())
		if errCallback != nil {
			errCallback()
		}
		return nil, err
	}
	lg.Info("vehicle created", "vehicleID", vehicle.ID)
	return createdVehicle, nil
}

func (s *VehicleService) UpdateVehicle(ctx context.Context, vehicle *entity.Vehicle, vehicleDTO dto.VehicleRequest) (*entity.Vehicle, error) {
	lg := s.logger.With("method", "UpdateVehicle")
	existingVehicle, err := s.vehicleRepo.GetByField(ctx, "id", vehicle.ID)
	if err != nil {
		if errors.Is(err, repository.ErrVehicleNotFound) {
			lg.Warn("vehicle not found for update", "vehicleID", vehicle.ID)
			return nil, repository.ErrVehicleNotFound
		}
		lg.Error("failed to get vehicle for update", "error", err.Error())
		return nil, err
	}

	var errCallback func()
	if vehicleDTO.ImageDocumentVehicle != nil {
		rs := bytes.NewReader(vehicleDTO.ImageDocumentVehicle)
		err := s.fr.UploadPublicFile(ctx, s.env.MinioProfilePictureBucket, vehicle.ID+"_image_document_vehicle", "application/octet-stream", rs)
		if err != nil {
			lg.Error("failed to upload document image", slog.Any("error", err))
			return nil, err
		}
		vehicle.ImageDocumentVehicle = fmt.Sprintf("%s/%s", s.env.MinioProfilePictureBucket, vehicle.ID+"_image_document_vehicle")

		errCallback = func() {
			if err1 := s.deleteFile(ctx, vehicle.ImageDocumentVehicle, lg); err1 != nil {
				slog.Error("failed to delete file", slog.Any("error", err1))
			}
		}
	}

	if vehicleDTO.ImageCardVehicle != nil {
		rs := bytes.NewReader(vehicleDTO.ImageCardVehicle)
		err := s.fr.UploadPublicFile(ctx, s.env.MinioProfilePictureBucket, vehicle.ID+"_image_card_vehicle", "application/octet-stream", rs)
		if err != nil {
			lg.Error("failed to upload card image", slog.Any("error", err))
			if errCallback != nil {
				errCallback()
			}
			return nil, err
		}
		vehicle.ImageCardVehicle = fmt.Sprintf("%s/%s", s.env.MinioProfilePictureBucket, vehicle.ID+"_image_card_vehicle")

		if errCallback != nil {
			oldCallback := errCallback
			errCallback = func() {
				oldCallback()
				if err1 := s.deleteFile(ctx, vehicle.ImageCardVehicle, lg); err1 != nil {
					slog.Error("failed to delete file", slog.Any("error", err1))
				}
			}
		} else {
			errCallback = func() {
				if err1 := s.deleteFile(ctx, vehicle.ImageCardVehicle, lg); err1 != nil {
					slog.Error("failed to delete file", slog.Any("error", err1))
				}
			}
		}
	}

	if vehicleDTO.ThirdPartyInsuranceImage != nil {
		rs := bytes.NewReader(vehicleDTO.ThirdPartyInsuranceImage)
		err := s.fr.UploadPublicFile(ctx, s.env.MinioProfilePictureBucket, vehicle.ID+"_third_party_insurance", "application/octet-stream", rs)
		if err != nil {
			lg.Error("failed to upload third party insurance image", slog.Any("error", err))
			if errCallback != nil {
				errCallback()
			}
			return nil, err
		}
		vehicle.ThirdPartyInsuranceImage = fmt.Sprintf("%s/%s", s.env.MinioProfilePictureBucket, vehicle.ID+"_third_party_insurance")

		if errCallback != nil {
			oldCallback := errCallback
			errCallback = func() {
				oldCallback()
				if err1 := s.deleteFile(ctx, vehicle.ThirdPartyInsuranceImage, lg); err1 != nil {
					slog.Error("failed to delete file", slog.Any("error", err1))
				}
			}
		} else {
			errCallback = func() {
				if err1 := s.deleteFile(ctx, vehicle.ThirdPartyInsuranceImage, lg); err1 != nil {
					slog.Error("failed to delete file", slog.Any("error", err1))
				}
			}
		}
	}

	if vehicleDTO.BodyInsuranceImage != nil {
		rs := bytes.NewReader(vehicleDTO.BodyInsuranceImage)
		err := s.fr.UploadPublicFile(ctx, s.env.MinioProfilePictureBucket, vehicle.ID+"_body_insurance", "application/octet-stream", rs)
		if err != nil {
			lg.Error("failed to upload body insurance image", slog.Any("error", err))
			if errCallback != nil {
				errCallback()
			}
			return nil, err
		}
		vehicle.BodyInsuranceImage = fmt.Sprintf("%s/%s", s.env.MinioProfilePictureBucket, vehicle.ID+"_body_insurance")

		if errCallback != nil {
			oldCallback := errCallback
			errCallback = func() {
				oldCallback()
				if err1 := s.deleteFile(ctx, vehicle.BodyInsuranceImage, lg); err1 != nil {
					slog.Error("failed to delete file", slog.Any("error", err1))
				}
			}
		} else {
			errCallback = func() {
				if err1 := s.deleteFile(ctx, vehicle.BodyInsuranceImage, lg); err1 != nil {
					slog.Error("failed to delete file", slog.Any("error", err1))
				}
			}
		}
	}

	lg.Info("updating vehicle", "vehicleID", existingVehicle.ID)
	_, err = s.vehicleRepo.Update(ctx, vehicle)
	if err != nil {
		lg.Error("failed to update vehicle", "error", err.Error())
		if errCallback != nil {
			errCallback()
		}
		return nil, err
	}
	updatedVehicle, err := s.vehicleRepo.GetByField(ctx, "id", existingVehicle.ID)
	if err != nil {
		lg.Error("failed to find updated vehicle", "error", err.Error())
		return nil, err
	}
	lg.Info("vehicle updated", "vehicleID", updatedVehicle.ID)
	return updatedVehicle, nil
}

func (s *VehicleService) deleteFile(ctx context.Context, file string, lg *slog.Logger) error {
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

func (s *VehicleService) ListVehicles(ctx context.Context, filters map[string]interface{}, sortField, sortOrder string, page, pageSize int) ([]entity.Vehicle, int64, error) {
	logger := s.logger.With("method", "ListVehicles")
	vehicles, total, err := s.vehicleRepo.List(ctx, filters, sortField, sortOrder, page, pageSize)
	if err != nil {
		logger.Error("failed to list vehicles", "error", err.Error())
		return nil, 0, err
	}
	logger.Info("vehicles listed", "totalVehicles", total)
	return vehicles, total, nil
}

func (s *VehicleService) DeleteVehicle(ctx context.Context, id string) error {
	logger := s.logger.With("method", "DeleteVehicle")
	existingVehicle, err := s.vehicleRepo.GetByField(ctx, "id", id)
	if err != nil {
		if errors.Is(err, repository.ErrVehicleNotFound) {
			logger.Warn("vehicle not found for deletion", "vehicleID", id)
			return nil
		}
		logger.Error("failed to find vehicle for deletion", "error", err.Error())
		return err
	}

	logger.Info("deleting vehicle", "vehicleID", existingVehicle.ID)
	err = s.vehicleRepo.Delete(ctx, id)
	if err != nil {
		logger.Error("failed to delete vehicle", "error", err.Error())
		return err
	}
	logger.Info("vehicle deleted", "vehicleID", existingVehicle.ID)
	return nil
}

func (s *VehicleService) Detail(ctx context.Context, id, value string) (*entity.Vehicle, error) {
	logger := s.logger.With("method", "Detail")
	vehicle, err := s.vehicleRepo.GetByField(ctx, id, value)
	if err != nil {
		if errors.Is(err, repository.ErrVehicleNotFound) {
			logger.Warn("vehicle not found for detail", "field", id, "value", value)
			return nil, errors.New("vehicle not found")
		}
		logger.Error("failed to get vehicle details", "error", err.Error())
		return nil, err
	}
	logger.Info("vehicle details retrieved", "vehicleID", vehicle.ID)
	return vehicle, nil
}
