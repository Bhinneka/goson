package goson

import (
	"errors"
	"reflect"
)

// this utility will hide json response based on struct json field you choose to hide

var errorsShouldBeStruct = errors.New("input should be struct")

type fieldSet map[string]bool

func genFieldSet(fields ...string) fieldSet {
	var fieldMap = make(fieldSet, len(fields))
	for _, v := range fields {
		fieldMap[v] = true
	}
	return fieldMap
}

// HideFields will hide json field based on struct json field you choose
// params
// - input (struct only)
// - fields json field that you will be hide
func HideFields(input interface{}, fields ...string) (map[string]interface{}, error) {
	typeInput := reflect.TypeOf(input)

	// validate input
	// if input is not struct return early errors
	if typeInput.Kind() != reflect.Struct {
		return nil, errorsShouldBeStruct
	}

	fieldSet := genFieldSet(fields...)

	output := make(map[string]interface{}, typeInput.NumField())

	valueInput := reflect.ValueOf(input)

	for i := 0; i < typeInput.NumField(); i++ {
		// field return each file on the struct
		field := typeInput.Field(i)

		// get json Tag
		// eg: Name `json:"name"` => result will be `name`
		jsonField := field.Tag.Get("json")
		if !fieldSet[jsonField] {
			output[jsonField] = valueInput.Field(i).Interface()
		}
	}
	return output, nil
}
