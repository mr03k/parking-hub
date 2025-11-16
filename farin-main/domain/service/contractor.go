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

type ContractorService struct {
	logger         *slog.Logger
	contractorRepo *repository.ContractorRepository
	fr             uploader.FileRepository
	env            *godotenv.Env
}

func NewContractorService(logger *slog.Logger, contractorRepo *repository.ContractorRepository, fr uploader.FileRepository,
	env *godotenv.Env) *ContractorService {
	return &ContractorService{
		logger:         logger.With("layer", "ContractorService"),
		contractorRepo: contractorRepo,
		fr:             fr,
		env:            env,
	}
}

func (s *ContractorService) CreateContractor(ctx context.Context, contractor *entity.Contractor) (*entity.Contractor, error) {
	logger := s.logger.With("method", "CreateContractor")
	contractor.ID = uuid.NewString()
	createdContractor, err := s.contractorRepo.Create(ctx, contractor)
	if err != nil {
		logger.Error("failed to create contractor", "error", err.Error())
		return nil, err
	}
	logger.Info("contractor created", "contractorID", contractor.ID)
	return createdContractor, nil
}

func (s *ContractorService) ListContractors(ctx context.Context, filters map[string]interface{}, sortField, sortOrder string, page, pageSize int) ([]entity.Contractor, int64, error) {
	logger := s.logger.With("method", "ListContractors")
	contractors, total, err := s.contractorRepo.List(ctx, filters, sortField, sortOrder, page, pageSize)
	if err != nil {
		logger.Error("failed to list contractors", "error", err.Error())
		return nil, 0, err
	}
	logger.Info("contractors listed", "totalContractors", total)
	return contractors, total, nil
}

func (s *ContractorService) UpdateContractor(ctx context.Context, contractor *entity.Contractor) (*entity.Contractor, error) {
	logger := s.logger.With("method", "UpdateContractor")
	existingContractor, err := s.contractorRepo.GetByField(ctx, "id", contractor.ID)
	if err != nil {
		if errors.Is(err, repository.ErrContractorNotFound) {
			logger.Warn("contractor not found for update", "contractorID", contractor.ID)
			return nil, repository.ErrContractorNotFound
		}
		logger.Error("failed to get contractor for update", "error", err.Error())
		return nil, err
	}

	logger.Info("updating contractor", "contractorID", existingContractor.ID)
	updatedContractor, err := s.contractorRepo.Update(ctx, contractor)
	if err != nil {
		logger.Error("failed to update contractor", "error", err.Error())
		return nil, err
	}
	logger.Info("contractor updated", "contractorID", updatedContractor.ID)
	return updatedContractor, nil
}

func (s *ContractorService) DeleteContractor(ctx context.Context, id string) error {
	logger := s.logger.With("method", "DeleteContractor")
	existingContractor, err := s.contractorRepo.GetByField(ctx, "id", id)
	if err != nil {
		if errors.Is(err, repository.ErrContractorNotFound) {
			logger.Warn("contractor not found for deletion", "contractorID", id)
			return nil
		}
		logger.Error("failed to find contractor for deletion", "error", err.Error())
		return err
	}

	logger.Info("deleting contractor", "contractorID", existingContractor.ID)
	err = s.contractorRepo.Delete(ctx, id)
	if err != nil {
		logger.Error("failed to delete contractor", "error", err.Error())
		return err
	}
	logger.Info("contractor deleted", "contractorID", existingContractor.ID)
	return nil
}

func (s *ContractorService) Detail(ctx context.Context, id, value string) (*entity.Contractor, error) {
	logger := s.logger.With("method", "Detail")
	contractor, err := s.contractorRepo.GetByField(ctx, id, value)
	if err != nil {
		if errors.Is(err, repository.ErrContractorNotFound) {
			logger.Warn("contractor not found for detail", "field", id, "value", value)
			return nil, errors.New("contractor not found")
		}
		logger.Error("failed to get contractor details", "error", err.Error())
		return nil, err
	}
	logger.Info("contractor details retrieved", "contractorID", contractor.ID)
	return contractor, nil
}
