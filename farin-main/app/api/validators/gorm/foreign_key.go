package gorm

import (
	gormdb "farin/infrastructure/gorm"
	"github.com/go-playground/validator/v10"
	"reflect"
)

type FkValidator struct {
	db *gormdb.GORMDB
}

func NewFkValidator(database *gormdb.GORMDB) FkValidator {
	return FkValidator{
		db: database,
	}
}

// fk validator
// please send destionation table name as foreign key
func (uv FkValidator) Handler() func(fl validator.FieldLevel) bool {
	return func(fl validator.FieldLevel) bool {
		field := fl.Field()
		table := fl.Param()
		var count int64

		// Get the appropriate value based on field type
		var value interface{}
		switch field.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			value = field.Int()
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			value = field.Uint()
		case reflect.String:
			value = field.String()
		case reflect.Interface:
			value = field.Interface()
		default:
			// Unsupported type
			return false
		}

		// Skip validation if value is zero/empty
		if value == nil || value == "" || value == 0 {
			return true
		}

		// Check if the ID exists in the referenced table
		uv.db.DB.Table(table).Where("id = ?", value).Count(&count)
		return count > 0
	}
}
