package service

import (
	"context"
	"errors"
	"farin/domain/entity"
	"farin/domain/repository"
	"farin/infrastructure/godotenv"
	"github.com/google/uuid"
	"github.com/mahdimehrabi/uploader"
	"log/slog"
	"time"
)

type DriverAssignmentService struct {
	logger               *slog.Logger
	driverAssignmentRepo *repository.DriverAssignmentRepository
	fr                   uploader.FileRepository
	env                  *godotenv.Env
}

func NewDriverAssignmentService(logger *slog.Logger, driverAssignmentRepo *repository.DriverAssignmentRepository, fr uploader.FileRepository,
	env *godotenv.Env) *DriverAssignmentService {
	return &DriverAssignmentService{
		logger:               logger.With("layer", "DriverAssignmentService"),
		driverAssignmentRepo: driverAssignmentRepo,
		fr:                   fr,
		env:                  env,
	}
}

func (s *DriverAssignmentService) CreateDriverAssignment(ctx context.Context, driverAssignment *entity.DriverAssignment) (*entity.DriverAssignment, error) {
	logger := s.logger.With("method", "CreateDriverAssignment")
	driverAssignment.ID = uuid.NewString()
	driverAssignment.CreatedAt = time.Now().Unix()
	driverAssignment.UpdatedAt = time.Now().Unix()
	createdDriverAssignment, err := s.driverAssignmentRepo.Create(ctx, driverAssignment)
	if err != nil {
		logger.Error("failed to create driverAssignment", "error", err.Error())
		return nil, err
	}
	logger.Info("driverAssignment created", "driverAssignmentID", driverAssignment.ID)
	return createdDriverAssignment, nil
}

func (s *DriverAssignmentService) ListDriverAssignments(ctx context.Context, filters map[string]interface{}, sortField, sortOrder string, page, pageSize int) ([]entity.DriverAssignment, int64, error) {
	logger := s.logger.With("method", "ListDriverAssignments")
	driverAssignments, total, err := s.driverAssignmentRepo.List(ctx, filters, sortField, sortOrder, page, pageSize)
	if err != nil {
		logger.Error("failed to list driverAssignments", "error", err.Error())
		return nil, 0, err
	}
	logger.Info("driverAssignments listed", "totalDriverAssignments", total)
	return driverAssignments, total, nil
}

func (s *DriverAssignmentService) UpdateDriverAssignment(ctx context.Context, driverAssignment *entity.DriverAssignment) (*entity.DriverAssignment, error) {
	logger := s.logger.With("method", "UpdateDriverAssignment")
	existingDriverAssignment, err := s.driverAssignmentRepo.GetByField(ctx, "id", driverAssignment.ID)
	if err != nil {
		if errors.Is(err, repository.ErrDriverAssignmentNotFound) {
			logger.Warn("driverAssignment not found for update", "driverAssignmentID", driverAssignment.ID)
			return nil, repository.ErrDriverAssignmentNotFound
		}
		logger.Error("failed to get driverAssignment for update", "error", err.Error())
		return nil, err
	}

	logger.Info("updating driverAssignment", "driverAssignmentID", existingDriverAssignment.ID)
	updatedDriverAssignment, err := s.driverAssignmentRepo.Update(ctx, driverAssignment)
	if err != nil {
		logger.Error("failed to update driverAssignment", "error", err.Error())
		return nil, err
	}
	logger.Info("driverAssignment updated", "driverAssignmentID", updatedDriverAssignment.ID)
	return updatedDriverAssignment, nil
}

func (s *DriverAssignmentService) DeleteDriverAssignment(ctx context.Context, id string) error {
	logger := s.logger.With("method", "DeleteDriverAssignment")
	existingDriverAssignment, err := s.driverAssignmentRepo.GetByField(ctx, "id", id)
	if err != nil {
		if errors.Is(err, repository.ErrDriverAssignmentNotFound) {
			logger.Warn("driverAssignment not found for deletion", "driverAssignmentID", id)
			return nil
		}
		logger.Error("failed to find driverAssignment for deletion", "error", err.Error())
		return err
	}

	logger.Info("deleting driverAssignment", "driverAssignmentID", existingDriverAssignment.ID)
	err = s.driverAssignmentRepo.Delete(ctx, id)
	if err != nil {
		logger.Error("failed to delete driverAssignment", "error", err.Error())
		return err
	}
	logger.Info("driverAssignment deleted", "driverAssignmentID", existingDriverAssignment.ID)
	return nil
}

func (s *DriverAssignmentService) Detail(ctx context.Context, id, value string) (*entity.DriverAssignment, error) {
	logger := s.logger.With("method", "Detail")
	driverAssignment, err := s.driverAssignmentRepo.GetByField(ctx, id, value)
	if err != nil {
		if errors.Is(err, repository.ErrDriverAssignmentNotFound) {
			logger.Warn("driverAssignment not found for detail", "field", id, "value", value)
			return nil, errors.New("driverAssignment not found")
		}
		logger.Error("failed to get driverAssignment details", "error", err.Error())
		return nil, err
	}
	logger.Info("driverAssignment details retrieved", "driverAssignmentID", driverAssignment.ID)
	return driverAssignment, nil
}
