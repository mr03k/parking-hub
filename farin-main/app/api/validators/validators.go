package validators

import (
	"farin/app/api/validators/gorm"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(gorm.NewFkValidator, NewTimestampValidator, gorm.NewUniqueValidator, NewFileValidator)

type Validators struct {
	uv   gorm.UniqueValidator
	fkv  gorm.FkValidator
	ts   TimestampValidator
	file FileValidator
}

func NewValidators(uv gorm.UniqueValidator, fkv gorm.FkValidator, ts TimestampValidator, file FileValidator) Validators {
	return Validators{
		uv:   uv,
		fkv:  fkv,
		ts:   ts,
		file: file,
	}
}

// Setup sets up middlewares
func (val Validators) Setup() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("uniqueGorm", val.uv.Handler())
		v.RegisterValidation("fkGorm", val.fkv.Handler())
		v.RegisterValidation("timestamp", val.ts.Handler())
		v.RegisterValidation("fileData", val.file.Handler())
	}
}
