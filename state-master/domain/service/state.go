package service

import (
	"context"
	"errors"
	"git.abanppc.com/farin-project/state/domain/entity"
	"git.abanppc.com/farin-project/state/domain/repository"
	"git.abanppc.com/farin-project/state/infrastructure/godotenv"
	"github.com/google/uuid"
	"log/slog"
)

type StateService struct {
	logger    *slog.Logger
	stateRepo *repository.StateRepository
	env       *godotenv.Env
}

func NewStateService(logger *slog.Logger, stateRepo *repository.StateRepository,
	env *godotenv.Env) *StateService {
	return &StateService{
		logger:    logger.With("layer", "StateService"),
		stateRepo: stateRepo,
		env:       env,
	}
}

func (s *StateService) CreateState(ctx context.Context, state *entity.State) (*entity.State, error) {
	logger := s.logger.With("method", "CreateState")
	state.ID = uuid.NewString()
	createdState, err := s.stateRepo.Create(ctx, state)
	if err != nil {
		logger.Error("failed to create state", "error", err.Error())
		return nil, err
	}
	logger.Info("state created", "stateID", state.ID)
	return createdState, nil
}

func (s *StateService) ListStates(ctx context.Context, filters map[string]interface{}, sortField, sortOrder string, page, pageSize int) ([]entity.State, int64, error) {
	logger := s.logger.With("method", "ListStates")
	states, total, err := s.stateRepo.List(ctx, filters, sortField, sortOrder, page, pageSize)
	if err != nil {
		logger.Error("failed to list states", "error", err.Error())
		return nil, 0, err
	}
	logger.Info("states listed", "totalStates", total)
	return states, total, nil
}

func (s *StateService) UpdateState(ctx context.Context, state *entity.State) (*entity.State, error) {
	logger := s.logger.With("method", "UpdateState")
	existingState, err := s.stateRepo.GetByField(ctx, "id", state.ID)
	if err != nil {
		if errors.Is(err, repository.ErrStateNotFound) {
			logger.Warn("state not found for update", "stateID", state.ID)
			return nil, repository.ErrStateNotFound
		}
		logger.Error("failed to get state for update", "error", err.Error())
		return nil, err
	}

	logger.Info("updating state", "stateID", existingState.ID)
	updatedState, err := s.stateRepo.Update(ctx, state)
	if err != nil {
		logger.Error("failed to update state", "error", err.Error())
		return nil, err
	}
	logger.Info("state updated", "stateID", updatedState.ID)
	return updatedState, nil
}

func (s *StateService) DeleteState(ctx context.Context, id string) error {
	logger := s.logger.With("method", "DeleteState")
	existingState, err := s.stateRepo.GetByField(ctx, "id", id)
	if err != nil {
		if errors.Is(err, repository.ErrStateNotFound) {
			logger.Warn("state not found for deletion", "stateID", id)
			return nil
		}
		logger.Error("failed to find state for deletion", "error", err.Error())
		return err
	}

	logger.Info("deleting state", "stateID", existingState.ID)
	err = s.stateRepo.Delete(ctx, id)
	if err != nil {
		logger.Error("failed to delete state", "error", err.Error())
		return err
	}
	logger.Info("state deleted", "stateID", existingState.ID)
	return nil
}

func (s *StateService) Detail(ctx context.Context, id, value string) (*entity.State, error) {
	logger := s.logger.With("method", "Detail")
	state, err := s.stateRepo.GetByField(ctx, id, value)
	if err != nil {
		if errors.Is(err, repository.ErrStateNotFound) {
			logger.Warn("state not found for detail", "field", id, "value", value)
			return nil, errors.New("state not found")
		}
		logger.Error("failed to get state details", "error", err.Error())
		return nil, err
	}
	logger.Info("state details retrieved", "stateID", state.ID)
	return state, nil
}
