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
)

type CalenderService struct {
	logger       *slog.Logger
	calenderRepo *repository.CalenderRepository
	fr           uploader.FileRepository
	env          *godotenv.Env
}

func NewCalenderService(logger *slog.Logger, calenderRepo *repository.CalenderRepository, fr uploader.FileRepository,
	env *godotenv.Env) *CalenderService {
	return &CalenderService{
		logger:       logger.With("layer", "CalenderService"),
		calenderRepo: calenderRepo,
		fr:           fr,
		env:          env,
	}
}

func (s *CalenderService) CreateCalender(ctx context.Context, calender *entity.Calender) (*entity.Calender, error) {
	logger := s.logger.With("method", "CreateCalender")
	calender.ID = uuid.NewString()
	createdCalender, err := s.calenderRepo.Create(ctx, calender)
	if err != nil {
		logger.Error("failed to create calender", "error", err.Error())
		return nil, err
	}
	logger.Info("calender created", "calenderID", calender.ID)
	return createdCalender, nil
}

func (s *CalenderService) ListCalenders(ctx context.Context, filters map[string]interface{}, sortField, sortOrder string, page, pageSize int) ([]entity.Calender, int64, error) {
	logger := s.logger.With("method", "ListCalenders")
	calenders, total, err := s.calenderRepo.List(ctx, filters, sortField, sortOrder, page, pageSize)
	if err != nil {
		logger.Error("failed to list calenders", "error", err.Error())
		return nil, 0, err
	}
	logger.Info("calenders listed", "totalCalenders", total)
	return calenders, total, nil
}

func (s *CalenderService) UpdateCalender(ctx context.Context, calender *entity.Calender) (*entity.Calender, error) {
	logger := s.logger.With("method", "UpdateCalender")
	existingCalender, err := s.calenderRepo.GetByField(ctx, "id", calender.ID)
	if err != nil {
		if errors.Is(err, repository.ErrCalenderNotFound) {
			logger.Warn("calender not found for update", "calenderID", calender.ID)
			return nil, repository.ErrCalenderNotFound
		}
		logger.Error("failed to get calender for update", "error", err.Error())
		return nil, err
	}

	logger.Info("updating calender", "calenderID", existingCalender.ID)
	updatedCalender, err := s.calenderRepo.Update(ctx, calender)
	if err != nil {
		logger.Error("failed to update calender", "error", err.Error())
		return nil, err
	}
	logger.Info("calender updated", "calenderID", updatedCalender.ID)
	return updatedCalender, nil
}

func (s *CalenderService) DeleteCalender(ctx context.Context, id string) error {
	logger := s.logger.With("method", "DeleteCalender")
	existingCalender, err := s.calenderRepo.GetByField(ctx, "id", id)
	if err != nil {
		if errors.Is(err, repository.ErrCalenderNotFound) {
			logger.Warn("calender not found for deletion", "calenderID", id)
			return nil
		}
		logger.Error("failed to find calender for deletion", "error", err.Error())
		return err
	}

	logger.Info("deleting calender", "calenderID", existingCalender.ID)
	err = s.calenderRepo.Delete(ctx, id)
	if err != nil {
		logger.Error("failed to delete calender", "error", err.Error())
		return err
	}
	logger.Info("calender deleted", "calenderID", existingCalender.ID)
	return nil
}

func (s *CalenderService) Detail(ctx context.Context, id, value string) (*entity.Calender, error) {
	logger := s.logger.With("method", "Detail")
	calender, err := s.calenderRepo.GetByField(ctx, id, value)
	if err != nil {
		if errors.Is(err, repository.ErrCalenderNotFound) {
			logger.Warn("calender not found for detail", "field", id, "value", value)
			return nil, errors.New("calender not found")
		}
		logger.Error("failed to get calender details", "error", err.Error())
		return nil, err
	}
	logger.Info("calender details retrieved", "calenderID", calender.ID)
	return calender, nil
}
