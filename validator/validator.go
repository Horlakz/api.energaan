package validator

import (
	"errors"
	"reflect"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/google/uuid"

	"github.com/horlakz/energaan-api/database"
)

type Validator[T any] struct{}

func (validator *Validator[T]) ValidateDBUnique(structure T, tableName string, uniqueField string, parameters map[string]interface{}) validation.RuleFunc {
	db := database.DatabaseFacade
	result := map[string]interface{}{}
	e := reflect.ValueOf(&structure).Elem()
	parentID := e.FieldByName("UUID").Interface().(uuid.UUID)

	return func(value interface{}) error {
		query := db.Table(tableName).Where(uniqueField+" = ?", value)

		if parentID != uuid.Nil {
			query = query.Where("uuid != ?", parentID.String())
		}

		for key, parameter := range parameters {
			param := e.FieldByName(key).Interface()
			query = query.Where(parameter.(string)+" = ?", param)
		}

		rows := query.Take(&result)

		if rows.RowsAffected > 0 {
			return errors.New("value already exist")
		}

		return nil
	}
}
