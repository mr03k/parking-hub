package gorm

import (
	gormdb "farin/infrastructure/gorm"
	"github.com/go-playground/validator/v10"
	"strings"
)

type UniqueValidator struct {
	database *gormdb.GORMDB
}

func NewUniqueValidator(database *gormdb.GORMDB) UniqueValidator {
	return UniqueValidator{
		database: database,
	}
}

// unique validator
// please send table name and column name as parmater
// splited by & like this uniqueDB=users&email
// for update please define a ID prop and fill that in your controller
func (uv UniqueValidator) Handler() func(fl validator.FieldLevel) bool {
	return func(fl validator.FieldLevel) bool {
		//extract id for check unique of update
		var ID string = ""
		IDField := fl.Parent().FieldByName("ID")
		if IDField.IsValid() && !IDField.IsZero() {
			ID = IDField.String()
		}

		value := fl.Field()
		str := strings.Split(fl.Param(), "&")
		table, column := str[0], str[1]
		var count int64
		if ID != "" {
			uv.database.DB.Table(table).Where("deleted_at=0").Where(column, value).Where("ID != ?", ID).Count(&count)
		} else {
			uv.database.DB.Table(table).Where("deleted_at=0").Where(column, value).Count(&count)
		}
		return count < 1
	}
}
