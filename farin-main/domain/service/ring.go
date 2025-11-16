package service

import (
	"context"
	"errors"
	"farin/domain/entity"
	"farin/domain/repository"
	"farin/infrastructure/godotenv"
	"github.com/mahdimehrabi/uploader"
	"log/slog"
)

type RingService struct {
	logger   *slog.Logger
	ringRepo *repository.RingRepository
	fr       uploader.FileRepository
	env      *godotenv.Env
}

func NewRingService(logger *slog.Logger, ringRepo *repository.RingRepository, fr uploader.FileRepository,
	env *godotenv.Env) *RingService {
	return &RingService{
		logger:   logger.With("layer", "RingService"),
		ringRepo: ringRepo,
		fr:       fr,
		env:      env,
	}
}

func (s *RingService) CreateRing(ctx context.Context, ring *entity.Ring) (*entity.Ring, error) {
	logger := s.logger.With("method", "CreateRing")
	createdRing, err := s.ringRepo.Create(ctx, ring)
	if err != nil {
		logger.Error("failed to create ring", "error", err.Error())
		return nil, err
	}
	logger.Info("ring created", "ringID", ring.ID)
	return createdRing, nil
}

func (s *RingService) ListRings(ctx context.Context, filters map[string]interface{}, sortField, sortOrder string, page, pageSize int) ([]entity.Ring, int64, error) {
	logger := s.logger.With("method", "ListRings")
	rings, total, err := s.ringRepo.List(ctx, filters, sortField, sortOrder, page, pageSize)
	if err != nil {
		logger.Error("failed to list rings", "error", err.Error())
		return nil, 0, err
	}
	logger.Info("rings listed", "totalRings", total)
	return rings, total, nil
}

func (s *RingService) UpdateRing(ctx context.Context, ring *entity.Ring) (*entity.Ring, error) {
	logger := s.logger.With("method", "UpdateRing")
	existingRing, err := s.ringRepo.GetByField(ctx, "id", ring.ID)
	if err != nil {
		if errors.Is(err, repository.ErrRingNotFound) {
			logger.Warn("ring not found for update", "ringID", ring.ID)
			return nil, repository.ErrRingNotFound
		}
		logger.Error("failed to get ring for update", "error", err.Error())
		return nil, err
	}

	logger.Info("updating ring", "ringID", existingRing.ID)
	updatedRing, err := s.ringRepo.Update(ctx, ring)
	if err != nil {
		logger.Error("failed to update ring", "error", err.Error())
		return nil, err
	}
	logger.Info("ring updated", "ringID", updatedRing.ID)
	return updatedRing, nil
}

func (s *RingService) DeleteRing(ctx context.Context, id string) error {
	logger := s.logger.With("method", "DeleteRing")
	existingRing, err := s.ringRepo.GetByField(ctx, "id", id)
	if err != nil {
		if errors.Is(err, repository.ErrRingNotFound) {
			logger.Warn("ring not found for deletion", "ringID", id)
			return nil
		}
		logger.Error("failed to find ring for deletion", "error", err.Error())
		return err
	}

	logger.Info("deleting ring", "ringID", existingRing.ID)
	err = s.ringRepo.Delete(ctx, id)
	if err != nil {
		logger.Error("failed to delete ring", "error", err.Error())
		return err
	}
	logger.Info("ring deleted", "ringID", existingRing.ID)
	return nil
}

func (s *RingService) Detail(ctx context.Context, id, value string) (*entity.Ring, error) {
	logger := s.logger.With("method", "Detail")
	ring, err := s.ringRepo.GetByField(ctx, id, value)
	if err != nil {
		if errors.Is(err, repository.ErrRingNotFound) {
			logger.Warn("ring not found for detail", "field", id, "value", value)
			return nil, errors.New("ring not found")
		}
		logger.Error("failed to get ring details", "error", err.Error())
		return nil, err
	}
	logger.Info("ring details retrieved", "ringID", ring.ID)
	return ring, nil
}
