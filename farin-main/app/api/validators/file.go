package validators

import (
	"bytes"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/mahdimehrabi/uploader"
	"log/slog"
	"strconv"
	"strings"
)

type FileValidator struct {
	logger *slog.Logger
}

func NewFileValidator(logger *slog.Logger) FileValidator {
	return FileValidator{
		logger: logger.With("layer", "FileValidator"),
	}
}

func (uv FileValidator) Handler() func(fl validator.FieldLevel) bool {
	return func(fl validator.FieldLevel) bool {
		value := fl.Field().Bytes()
		if len(value) < 1 {
			return true
		}
		strs := strings.Split(fl.Param(), ";")
		if len(strs) != 2 {
			return true
		}
		rs := bytes.NewReader(value)
		fileTypesStr := strs[0]
		if fileTypesStr != "" {
			fileTypes := strings.Split(fileTypesStr, "&")
			if len(fileTypes) < 1 {
				return true
			}
			_, err := uploader.ValidateFileType(rs, fileTypes...)
			if err != nil {
				if errors.Is(err, uploader.ErrInvalidFileType) {
					return false
				}
				uv.logger.Error("failed to validate file type", slog.With("type", fileTypesStr),
					slog.With("error", err.Error()))
				return false
			}
		}

		fileSize, err := strconv.Atoi(strs[1])
		if err != nil {
			return true
		}

		if err := uploader.ValidateFileSize(rs, int64(fileSize)); err != nil {
			if errors.Is(err, uploader.ErrFileTooLarge) {
				return false
			}
			uv.logger.Error("failed to validate file type", slog.With("fileSize", fileSize),
				slog.With("error", err.Error()))
			return false
		}
		return true
	}
}
