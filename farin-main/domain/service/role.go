package service

import (
	"context"
	"errors"
	"farin/domain/dto"
	"farin/domain/entity"
	"farin/domain/repository"
	"farin/infrastructure/godotenv"
	"github.com/google/uuid"
	"github.com/mahdimehrabi/uploader"
	"log/slog"
)

var ErrModifingSystemRole = errors.New("necessary system role cannot be deleted/modified")

type RoleService struct {
	logger   *slog.Logger
	RoleRepo *repository.RoleRepository
	fr       uploader.FileRepository
	env      *godotenv.Env
}

func NewRoleService(logger *slog.Logger, RoleRepo *repository.RoleRepository, fr uploader.FileRepository,
	env *godotenv.Env) *RoleService {
	return &RoleService{
		logger:   logger.With("layer", "RoleService"),
		RoleRepo: RoleRepo,
		fr:       fr,
		env:      env,
	}
}
func (s *RoleService) CreateRole(ctx context.Context, Role *entity.Role, RoleDTO *dto.RoleRequest) (*entity.Role, error) {
	lg := s.logger.With("method", "CreateRole")
	Role.ID = uuid.NewString()

	createdRole, err := s.RoleRepo.Create(ctx, Role)
	if err != nil {
		lg.Error("failed to create Role", "error", err.Error())
		return nil, err
	}
	lg.Info("Role created", "RoleID", Role.ID)
	return createdRole, nil
}

func (s *RoleService) ListRoles(ctx context.Context, filters map[string]interface{}, sortField, sortOrder string, page, pageSize int) ([]entity.Role, int64, error) {
	logger := s.logger.With("method", "ListRoles")
	Roles, total, err := s.RoleRepo.List(ctx, filters, sortField, sortOrder, page, pageSize)
	if err != nil {
		logger.Error("failed to list Roles", "error", err.Error())
		return nil, 0, err
	}
	logger.Info("Roles listed", "totalRoles", total)
	return Roles, total, nil
}

func (s *RoleService) UpdateRole(ctx context.Context, Role *entity.Role, RoleDTO dto.RoleRequest) (*entity.Role, error) {
	lg := s.logger.With("method", "UpdateRole")
	existingRole, err := s.RoleRepo.GetByField(ctx, "id", Role.ID)
	if err != nil {
		if errors.Is(err, repository.ErrRoleNotFound) {
			lg.Warn("Role not found for update", "RoleID", Role.ID)
			return nil, repository.ErrRoleNotFound
		}
		lg.Error("failed to get Role for update", "error", err.Error())
		return nil, err
	}
	if existingRole.Title == "Admin" || existingRole.Title == "Driver" {
		return nil, ErrModifingSystemRole
	}

	lg.Info("updating Role", "RoleID", existingRole.ID)
	_, err = s.RoleRepo.Update(ctx, Role)
	if err != nil {
		lg.Error("failed to update Role", "error", err.Error())
		return nil, err
	}
	updatedRole, err := s.RoleRepo.GetByField(ctx, "id", existingRole.ID)
	if err != nil {
		lg.Error("failed to  findupdate Role", "error", err.Error())
		return nil, err
	}
	lg.Info("Role updated", "RoleID", updatedRole.ID)
	return updatedRole, nil
}

func (s *RoleService) DeleteRole(ctx context.Context, id string) error {
	logger := s.logger.With("method", "DeleteRole")
	existingRole, err := s.RoleRepo.GetByField(ctx, "id", id)
	if err != nil {
		if errors.Is(err, repository.ErrRoleNotFound) {
			logger.Warn("Role not found for deletion", "RoleID", id)
			return nil
		}
		logger.Error("failed to find Role for deletion", "error", err.Error())
		return err
	}

	if existingRole.Title == "Admin" || existingRole.Title == "Driver" {
		return ErrModifingSystemRole
	}

	logger.Info("deleting Role", "RoleID", existingRole.ID)
	err = s.RoleRepo.Delete(ctx, id)
	if err != nil {
		logger.Error("failed to delete Role", "error", err.Error())
		return err
	}
	logger.Info("Role deleted", "RoleID", existingRole.ID)
	return nil
}

func (s *RoleService) Detail(ctx context.Context, id, value string) (*entity.Role, error) {
	logger := s.logger.With("method", "Detail")
	Role, err := s.RoleRepo.GetByField(ctx, id, value)
	if err != nil {
		if errors.Is(err, repository.ErrRoleNotFound) {
			logger.Warn("Role not found for detail", "field", id, "value", value)
			return nil, errors.New("Role not found")
		}
		logger.Error("failed to get Role details", "error", err.Error())
		return nil, err
	}
	logger.Info("Role details retrieved", "RoleID", Role.ID)
	return Role, nil
}
