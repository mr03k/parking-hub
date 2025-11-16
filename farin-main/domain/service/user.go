package service

import (
	"bytes"
	"context"
	"errors"
	"farin/domain/dto"
	"farin/domain/entity"
	"farin/domain/repository"
	"farin/infrastructure/godotenv"
	"farin/util/encrypt"
	"fmt"
	"github.com/google/uuid"
	"github.com/mahdimehrabi/uploader"
	"io"
	"log/slog"
	"strings"
)

var (
	ErrFileTooLarge = errors.New("file too large(max file size is 4096)")
)

type UserService struct {
	logger   *slog.Logger
	userRepo *repository.UserRepository
	fr       uploader.FileRepository
	env      *godotenv.Env
}

func NewUserService(logger *slog.Logger, userRepo *repository.UserRepository, fr uploader.FileRepository,
	env *godotenv.Env) *UserService {
	return &UserService{
		logger:   logger.With("layer", "UserService"),
		userRepo: userRepo,
		fr:       fr,
		env:      env,
	}
}

func (s *UserService) CreateUser(ctx context.Context, user *entity.User, userDTO dto.UserRequest) (*entity.User, error) {
	lg := s.logger.With("method", "CreateUser")
	user.ID = uuid.NewString()
	user.Password = encrypt.HashSHA256(user.Password)
	var errCallback func()
	if userDTO.ProfileImage != nil {
		rs := bytes.NewReader(userDTO.ProfileImage)
		err := s.fr.UploadPublicFile(ctx, s.env.MinioProfilePictureBucket, user.ID+"_profile_image", "application/octet-stream", rs)
		if err != nil {
			lg.Error("failed to upload picture", slog.Any("error", err))
			return nil, err
		}
		user.ProfileImage = fmt.Sprintf("%s/%s", s.env.MinioProfilePictureBucket, user.ID+"_profile_image")

		errCallback = func() {
			if err1 := s.deleteFile(ctx, user.ProfileImage, lg); err1 != nil {
				slog.Error("failed to delete file", slog.Any("error", err1))
			}
		}
	}
	createdUser, err := s.userRepo.Create(ctx, user)
	if err != nil {
		lg.Error("failed to create user", "error", err.Error())
		errCallback()
		return nil, err
	}
	lg.Info("user created", "userID", user.ID)
	return createdUser, nil
}

func (s *UserService) ListUsers(ctx context.Context, filters map[string]interface{}, sortField, sortOrder string, page, pageSize int) ([]entity.User, int64, error) {
	logger := s.logger.With("method", "ListUsers")
	users, total, err := s.userRepo.List(ctx, filters, sortField, sortOrder, page, pageSize)
	if err != nil {
		logger.Error("failed to list users", "error", err.Error())
		return nil, 0, err
	}
	logger.Info("users listed", "totalUsers", total)
	return users, total, nil
}

func (s *UserService) UpdateUser(ctx context.Context, user *entity.User, userDTO dto.UserRequest) (*entity.User, error) {
	lg := s.logger.With("method", "UpdateUser")
	existingUser, err := s.userRepo.GetByField(ctx, "id", user.ID)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			lg.Warn("user not found for update", "userID", user.ID)
			return nil, repository.ErrUserNotFound
		}
		lg.Error("failed to get user for update", "error", err.Error())
		return nil, err
	}
	var errCallback func()
	if userDTO.ProfileImage != nil {
		rs := bytes.NewReader(userDTO.ProfileImage)
		err := s.fr.UploadPublicFile(ctx, s.env.MinioProfilePictureBucket, user.ID+"_profile_image", "application/octet-stream", rs)
		if err != nil {
			lg.Error("failed to upload picture", slog.Any("error", err))
			return nil, err
		}
		user.ProfileImage = fmt.Sprintf("%s/%s", s.env.MinioProfilePictureBucket, user.ID+"_profile_image")

		errCallback = func() {
			if err1 := s.deleteFile(ctx, user.ProfileImage, lg); err1 != nil {
				slog.Error("failed to delete file", slog.Any("error", err1))
			}
		}
	}
	lg.Info("updating user", "userID", existingUser.ID)
	updatedUser, err := s.userRepo.Update(ctx, user)
	if err != nil {
		lg.Error("failed to update user", "error", err.Error())
		errCallback()
		return nil, err
	}
	lg.Info("user updated", "userID", updatedUser.ID)
	return updatedUser, nil
}

func (s *UserService) DeleteUser(ctx context.Context, id string) error {
	logger := s.logger.With("method", "DeleteUser")
	existingUser, err := s.userRepo.GetByField(ctx, "id", id)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			logger.Warn("user not found for deletion", "userID", id)
			return nil
		}
		logger.Error("failed to find user for deletion", "error", err.Error())
		return err
	}

	logger.Info("deleting user", "userID", existingUser.ID)
	err = s.userRepo.Delete(ctx, id)
	if err != nil {
		logger.Error("failed to delete user", "error", err.Error())
		return err
	}
	logger.Info("user deleted", "userID", existingUser.ID)
	return nil
}

func (s *UserService) Detail(ctx context.Context, id, value string) (*entity.User, error) {
	logger := s.logger.With("method", "Detail")
	user, err := s.userRepo.GetByField(ctx, id, value)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			logger.Warn("user not found for detail", "field", id, "value", value)
			return nil, errors.New("user not found")
		}
		logger.Error("failed to get user details", "error", err.Error())
		return nil, err
	}
	logger.Info("user details retrieved", "userID", user.ID)
	return user, nil
}

func (s *UserService) UploadPicture(ctx context.Context, user *entity.User, picture io.ReadSeeker) error {
	lg := s.logger.With("method", "UploadPicture")

	contentType, err := uploader.ValidateFileType(picture, "image/jpeg", "image/png")
	if err != nil {
		if errors.Is(err, uploader.ErrInvalidFileType) {
			return err
		}
		lg.Error("invalid file type", slog.Any("error", err))
		return err
	}

	if err := uploader.ValidateFileSize(picture, 256000); err != nil {
		if errors.Is(err, uploader.ErrFileTooLarge) {
			return ErrFileTooLarge
		}
		return err
	}

	err = s.fr.UploadPublicFile(ctx, s.env.MinioProfilePictureBucket, user.ID, contentType, picture)
	if err != nil {
		lg.Error("failed to upload picture", slog.Any("error", err))
		return err
	}
	errCallback := func() {
		if err1 := s.deleteFile(ctx, user.ProfileImage, lg); err1 != nil {
			slog.Error("failed to delete file", slog.Any("error", err1))
		}
	}
	user.ProfileImage = fmt.Sprintf("%s/%s", s.env.MinioProfilePictureBucket, user.ID)

	// Save the picture URL to the user's profile
	if _, err := s.userRepo.Update(ctx, user); err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			errCallback()
			return nil
		}
		errCallback()
		return err
	}

	return nil
}

func (s UserService) deleteFile(ctx context.Context, file string, lg *slog.Logger) error {
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
