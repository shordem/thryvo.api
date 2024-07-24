package validator

import (
	"encoding/json"
	"errors"
	"log"
	"reflect"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/google/uuid"

	"github.com/shordem/api.thryvo/lib/database"
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

func (validator *Validator[T]) ValidateErr(err error) (map[string]interface{}, error) {
	if e, ok := err.(validation.InternalError); ok {
		log.Println(e.InternalError())
		return nil, nil
	}

	var dat map[string]interface{}
	m, _ := json.Marshal(err)

	json.Unmarshal(m, &dat)
	return dat, err
}
