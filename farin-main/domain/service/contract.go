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

type ContractService struct {
	logger       *slog.Logger
	contractRepo *repository.ContractRepository
	fr           uploader.FileRepository
	env          *godotenv.Env
}

func NewContractService(logger *slog.Logger, contractRepo *repository.ContractRepository, fr uploader.FileRepository,
	env *godotenv.Env) *ContractService {
	return &ContractService{
		logger:       logger.With("layer", "ContractService"),
		contractRepo: contractRepo,
		fr:           fr,
		env:          env,
	}
}

func (s *ContractService) CreateContract(ctx context.Context, contract *entity.Contract) (*entity.Contract, error) {
	logger := s.logger.With("method", "CreateContract")
	contract.ID = uuid.NewString()
	createdContract, err := s.contractRepo.Create(ctx, contract)
	if err != nil {
		logger.Error("failed to create contract", "error", err.Error())
		return nil, err
	}
	logger.Info("contract created", "contractID", contract.ID)
	return createdContract, nil
}

func (s *ContractService) ListContracts(ctx context.Context, filters map[string]interface{}, sortField, sortOrder string, page, pageSize int) ([]entity.Contract, int64, error) {
	logger := s.logger.With("method", "ListContracts")
	contracts, total, err := s.contractRepo.List(ctx, filters, sortField, sortOrder, page, pageSize)
	if err != nil {
		logger.Error("failed to list contracts", "error", err.Error())
		return nil, 0, err
	}
	logger.Info("contracts listed", "totalContracts", total)
	return contracts, total, nil
}

func (s *ContractService) UpdateContract(ctx context.Context, contract *entity.Contract) (*entity.Contract, error) {
	logger := s.logger.With("method", "UpdateContract")
	existingContract, err := s.contractRepo.GetByField(ctx, "id", contract.ID)
	if err != nil {
		if errors.Is(err, repository.ErrContractNotFound) {
			logger.Warn("contract not found for update", "contractID", contract.ID)
			return nil, repository.ErrContractNotFound
		}
		logger.Error("failed to get contract for update", "error", err.Error())
		return nil, err
	}

	logger.Info("updating contract", "contractID", existingContract.ID)
	updatedContract, err := s.contractRepo.Update(ctx, contract)
	if err != nil {
		logger.Error("failed to update contract", "error", err.Error())
		return nil, err
	}
	logger.Info("contract updated", "contractID", updatedContract.ID)
	return updatedContract, nil
}

func (s *ContractService) DeleteContract(ctx context.Context, id string) error {
	logger := s.logger.With("method", "DeleteContract")
	existingContract, err := s.contractRepo.GetByField(ctx, "id", id)
	if err != nil {
		if errors.Is(err, repository.ErrContractNotFound) {
			logger.Warn("contract not found for deletion", "contractID", id)
			return nil
		}
		logger.Error("failed to find contract for deletion", "error", err.Error())
		return err
	}

	logger.Info("deleting contract", "contractID", existingContract.ID)
	err = s.contractRepo.Delete(ctx, id)
	if err != nil {
		logger.Error("failed to delete contract", "error", err.Error())
		return err
	}
	logger.Info("contract deleted", "contractID", existingContract.ID)
	return nil
}

func (s *ContractService) Detail(ctx context.Context, id, value string) (*entity.Contract, error) {
	logger := s.logger.With("method", "Detail")
	contract, err := s.contractRepo.GetByField(ctx, id, value)
	if err != nil {
		if errors.Is(err, repository.ErrContractNotFound) {
			logger.Warn("contract not found for detail", "field", id, "value", value)
			return nil, errors.New("contract not found")
		}
		logger.Error("failed to get contract details", "error", err.Error())
		return nil, err
	}
	logger.Info("contract details retrieved", "contractID", contract.ID)
	return contract, nil
}
