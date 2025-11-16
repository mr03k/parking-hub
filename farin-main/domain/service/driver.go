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

type DriverService struct {
	logger     *slog.Logger
	driverRepo *repository.DriverRepository
	fr         uploader.FileRepository
	env        *godotenv.Env
}

func NewDriverService(logger *slog.Logger, driverRepo *repository.DriverRepository, fr uploader.FileRepository,
	env *godotenv.Env) *DriverService {
	return &DriverService{
		logger:     logger.With("layer", "DriverService"),
		driverRepo: driverRepo,
		fr:         fr,
		env:        env,
	}
}

func (s *DriverService) CreateDriver(ctx context.Context, driver *entity.Driver, request dto.DriverRequest) (*entity.Driver, error) {
	lg := s.logger.With("method", "CreateDriver")
	driver.ID = uuid.NewString()

	var errCallback func()
	if request.DriverPhoto != nil {
		rs := bytes.NewReader(request.DriverPhoto)
		err := s.fr.UploadPublicFile(ctx, s.env.MinioProfilePictureBucket, driver.ID+"_driver_photo", "application/octet-stream", rs)
		if err != nil {
			lg.Error("failed to upload picture", slog.Any("error", err))
			return nil, err
		}
		driver.DriverPhoto = fmt.Sprintf("%s/%s", s.env.MinioProfilePictureBucket, driver.ID+"_driver_photo")

		errCallback = func() {
			errCallback()
			if err1 := s.deleteFile(ctx, driver.ID+"_driver_photo", lg); err1 != nil {
				slog.Error("failed to delete file", slog.Any("error", err1))
			}
		}
	}

	if request.BirthCertificateImage != nil {
		rs := bytes.NewReader(request.BirthCertificateImage)
		err := s.fr.UploadPublicFile(ctx, s.env.MinioProfilePictureBucket, driver.ID+"_birth_certificate_image", "application/octet-stream", rs)
		if err != nil {
			lg.Error("failed to upload picture", slog.Any("error", err))
			return nil, err
		}
		driver.BirthCertificateImage = fmt.Sprintf("%s/%s", s.env.MinioProfilePictureBucket, driver.ID+"_birth_certificate_image")

		errCallback = func() {
			errCallback()
			if err1 := s.deleteFile(ctx, driver.ID+"_birth_certificate_image", lg); err1 != nil {
				slog.Error("failed to delete file", slog.Any("error", err1))
			}
		}
	}

	if request.CriminalRecordImage != nil {
		rs := bytes.NewReader(request.CriminalRecordImage)
		err := s.fr.UploadPublicFile(ctx, s.env.MinioProfilePictureBucket, driver.ID+"_criminal_record_image", "application/octet-stream", rs)
		if err != nil {
			lg.Error("failed to upload picture", slog.Any("error", err))
			return nil, err
		}
		driver.CriminalRecordImage = fmt.Sprintf("%s/%s", s.env.MinioProfilePictureBucket, driver.ID+"_criminal_record_image")

		errCallback = func() {
			errCallback()
			if err1 := s.deleteFile(ctx, driver.ID+"_criminal_record_image", lg); err1 != nil {
				slog.Error("failed to delete file", slog.Any("error", err1))
			}
		}
	}

	if request.IDCardImage != nil {
		rs := bytes.NewReader(request.IDCardImage)
		err := s.fr.UploadPublicFile(ctx, s.env.MinioProfilePictureBucket, driver.ID+"_id_card_image", "application/octet-stream", rs)
		if err != nil {
			lg.Error("failed to upload picture", slog.Any("error", err))
			return nil, err
		}
		driver.IDCardImage = fmt.Sprintf("%s/%s", s.env.MinioProfilePictureBucket, driver.ID+"_id_card_image")

		errCallback = func() {
			errCallback()
			if err1 := s.deleteFile(ctx, driver.ID+"_id_card_image", lg); err1 != nil {
				slog.Error("failed to delete file", slog.Any("error", err1))
			}
		}
	}

	if request.HealthCertificateImage != nil {
		rs := bytes.NewReader(request.HealthCertificateImage)
		err := s.fr.UploadPublicFile(ctx, s.env.MinioProfilePictureBucket, driver.ID+"_health_certificate_image", "application/octet-stream", rs)
		if err != nil {
			lg.Error("failed to upload picture", slog.Any("error", err))
			return nil, err
		}
		driver.HealthCertificateImage = fmt.Sprintf("%s/%s", s.env.MinioProfilePictureBucket, driver.ID+"_health_certificate_image")

		errCallback = func() {
			errCallback()
			if err1 := s.deleteFile(ctx, driver.ID+"_health_certificate_image", lg); err1 != nil {
				slog.Error("failed to delete file", slog.Any("error", err1))
			}
		}
	}

	if request.MilitaryServiceCardImage != nil {
		rs := bytes.NewReader(request.MilitaryServiceCardImage)
		err := s.fr.UploadPublicFile(ctx, s.env.MinioProfilePictureBucket, driver.ID+"_military_service_image", "application/octet-stream", rs)
		if err != nil {
			lg.Error("failed to upload picture", slog.Any("error", err))
			return nil, err
		}
		driver.MilitaryServiceCardImage = fmt.Sprintf("%s/%s", s.env.MinioProfilePictureBucket, driver.ID+"_military_service_image")

		errCallback = func() {
			errCallback()
			if err1 := s.deleteFile(ctx, driver.ID+"_military_service_image", lg); err1 != nil {
				slog.Error("failed to delete file", slog.Any("error", err1))
			}
		}
	}

	createdDriver, err := s.driverRepo.Create(ctx, driver)
	if err != nil {
		lg.Error("failed to create driver", "error", err.Error())
		errCallback()
		return nil, err
	}
	lg.Info("driver created", "driverID", driver.ID)
	return createdDriver, nil
}

func (s DriverService) deleteFile(ctx context.Context, file string, lg *slog.Logger) error {
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

func (s *DriverService) ListDrivers(ctx context.Context, filters map[string]interface{}, sortField, sortOrder string, page, pageSize int) ([]entity.Driver, int64, error) {
	logger := s.logger.With("method", "ListDrivers")
	drivers, total, err := s.driverRepo.List(ctx, filters, sortField, sortOrder, page, pageSize)
	if err != nil {
		logger.Error("failed to list drivers", "error", err.Error())
		return nil, 0, err
	}
	logger.Info("drivers listed", "totalDrivers", total)
	return drivers, total, nil
}

func (s *DriverService) UpdateDriver(ctx context.Context, driver *entity.Driver, request dto.DriverRequest) (*entity.Driver, error) {
	lg := s.logger.With("method", "UpdateDriver")
	existingDriver, err := s.driverRepo.GetByField(ctx, "id", driver.ID)
	if err != nil {
		if errors.Is(err, repository.ErrDriverNotFound) {
			lg.Warn("driver not found for update", "driverID", driver.ID)
			return nil, repository.ErrDriverNotFound
		}
		lg.Error("failed to get driver for update", "error", err.Error())
		return nil, err
	}

	var errCallback func()
	if request.DriverPhoto != nil {
		rs := bytes.NewReader(request.DriverPhoto)
		err := s.fr.UploadPublicFile(ctx, s.env.MinioProfilePictureBucket, driver.ID+"_driver_photo", "application/octet-stream", rs)
		if err != nil {
			lg.Error("failed to upload picture", slog.Any("error", err))
			return nil, err
		}
		driver.DriverPhoto = fmt.Sprintf("%s/%s", s.env.MinioProfilePictureBucket, driver.ID+"_driver_photo")

		errCallback = func() {
			errCallback()
			if err1 := s.deleteFile(ctx, driver.ID+"_driver_photo", lg); err1 != nil {
				slog.Error("failed to delete file", slog.Any("error", err1))
			}
		}
	}

	if request.BirthCertificateImage != nil {
		rs := bytes.NewReader(request.BirthCertificateImage)
		err := s.fr.UploadPublicFile(ctx, s.env.MinioProfilePictureBucket, driver.ID+"_birth_certificate_image", "application/octet-stream", rs)
		if err != nil {
			lg.Error("failed to upload picture", slog.Any("error", err))
			return nil, err
		}
		driver.BirthCertificateImage = fmt.Sprintf("%s/%s", s.env.MinioProfilePictureBucket, driver.ID+"_birth_certificate_image")

		errCallback = func() {
			errCallback()
			if err1 := s.deleteFile(ctx, driver.ID+"_birth_certificate_image", lg); err1 != nil {
				slog.Error("failed to delete file", slog.Any("error", err1))
			}
		}
	}

	if request.CriminalRecordImage != nil {
		rs := bytes.NewReader(request.CriminalRecordImage)
		err := s.fr.UploadPublicFile(ctx, s.env.MinioProfilePictureBucket, driver.ID+"_criminal_record_image", "application/octet-stream", rs)
		if err != nil {
			lg.Error("failed to upload picture", slog.Any("error", err))
			return nil, err
		}
		driver.CriminalRecordImage = fmt.Sprintf("%s/%s", s.env.MinioProfilePictureBucket, driver.ID+"_criminal_record_image")

		errCallback = func() {
			errCallback()
			if err1 := s.deleteFile(ctx, driver.ID+"_criminal_record_image", lg); err1 != nil {
				slog.Error("failed to delete file", slog.Any("error", err1))
			}
		}
	}

	if request.IDCardImage != nil {
		rs := bytes.NewReader(request.IDCardImage)
		err := s.fr.UploadPublicFile(ctx, s.env.MinioProfilePictureBucket, driver.ID+"_id_card_image", "application/octet-stream", rs)
		if err != nil {
			lg.Error("failed to upload picture", slog.Any("error", err))
			return nil, err
		}
		driver.IDCardImage = fmt.Sprintf("%s/%s", s.env.MinioProfilePictureBucket, driver.ID+"_id_card_image")

		errCallback = func() {
			errCallback()
			if err1 := s.deleteFile(ctx, driver.ID+"_id_card_image", lg); err1 != nil {
				slog.Error("failed to delete file", slog.Any("error", err1))
			}
		}
	}

	if request.HealthCertificateImage != nil {
		rs := bytes.NewReader(request.HealthCertificateImage)
		err := s.fr.UploadPublicFile(ctx, s.env.MinioProfilePictureBucket, driver.ID+"_health_certificate_image", "application/octet-stream", rs)
		if err != nil {
			lg.Error("failed to upload picture", slog.Any("error", err))
			return nil, err
		}
		driver.HealthCertificateImage = fmt.Sprintf("%s/%s", s.env.MinioProfilePictureBucket, driver.ID+"_health_certificate_image")

		errCallback = func() {
			errCallback()
			if err1 := s.deleteFile(ctx, driver.ID+"_health_certificate_image", lg); err1 != nil {
				slog.Error("failed to delete file", slog.Any("error", err1))
			}
		}
	}

	if request.MilitaryServiceCardImage != nil {
		rs := bytes.NewReader(request.MilitaryServiceCardImage)
		err := s.fr.UploadPublicFile(ctx, s.env.MinioProfilePictureBucket, driver.ID+"_military_service_image", "application/octet-stream", rs)
		if err != nil {
			lg.Error("failed to upload picture", slog.Any("error", err))
			return nil, err
		}
		driver.MilitaryServiceCardImage = fmt.Sprintf("%s/%s", s.env.MinioProfilePictureBucket, driver.ID+"_military_service_image")

		errCallback = func() {
			errCallback()
			if err1 := s.deleteFile(ctx, driver.ID+"_military_service_image", lg); err1 != nil {
				slog.Error("failed to delete file", slog.Any("error", err1))
			}
		}
	}

	lg.Info("updating driver", "driverID", existingDriver.ID)
	updatedDriver, err := s.driverRepo.Update(ctx, driver)
	if err != nil {
		lg.Error("failed to update driver", "error", err.Error())
		errCallback()
		return nil, err
	}
	lg.Info("driver updated", "driverID", updatedDriver.ID)
	return updatedDriver, nil
}

func (s *DriverService) DeleteDriver(ctx context.Context, id string) error {
	logger := s.logger.With("method", "DeleteDriver")
	existingDriver, err := s.driverRepo.GetByField(ctx, "id", id)
	if err != nil {
		if errors.Is(err, repository.ErrDriverNotFound) {
			logger.Warn("driver not found for deletion", "driverID", id)
			return nil
		}
		logger.Error("failed to find driver for deletion", "error", err.Error())
		return err
	}

	logger.Info("deleting driver", "driverID", existingDriver.ID)
	err = s.driverRepo.Delete(ctx, id)
	if err != nil {
		logger.Error("failed to delete driver", "error", err.Error())
		return err
	}
	logger.Info("driver deleted", "driverID", existingDriver.ID)
	return nil
}

func (s *DriverService) Detail(ctx context.Context, id, value string) (*entity.Driver, error) {
	logger := s.logger.With("method", "Detail")
	driver, err := s.driverRepo.GetByField(ctx, id, value)
	if err != nil {
		if errors.Is(err, repository.ErrDriverNotFound) {
			logger.Warn("driver not found for detail", "field", id, "value", value)
			return nil, errors.New("driver not found")
		}
		logger.Error("failed to get driver details", "error", err.Error())
		return nil, err
	}
	logger.Info("driver details retrieved", "driverID", driver.ID)
	return driver, nil
}
