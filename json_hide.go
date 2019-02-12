package goson

import (
	"errors"
	"reflect"
)

// this utility will hide json response based on struct json field you choose to hide

var errorsShouldBeStruct = errors.New("input should be struct")

type fieldSet map[string]bool

// HideFields will hide json field based on struct json field you choose
// params
// - input (struct only)
// - fields json field that you will be hide
func HideFields(input interface{}, fields ...string) (map[string]interface{}, error) {
	typeInput := reflect.TypeOf(input)
	valueInput := reflect.ValueOf(input)

	// validate input
	// if input is not struct return early errors
	if typeInput.Kind() != reflect.Struct {
		return nil, errorsShouldBeStruct
	}

	fs := genFieldSet(fields...)

	return hideFields(typeInput, valueInput, fs), nil
}

func genFieldSet(fields ...string) fieldSet {
	var fieldMap = make(fieldSet, len(fields))
	for _, v := range fields {
		fieldMap[v] = true
	}
	return fieldMap
}

func hideFields(t reflect.Type, v reflect.Value, fs fieldSet) map[string]interface{} {
	output := make(map[string]interface{}, t.NumField())

	for i := 0; i < t.NumField(); i++ {
		// field return each file on the struct
		field := v.Field(i)
		fieldType := t.Field(i)

		if field.Kind() == reflect.Interface && !v.IsNil() {
			elm := field.Elem()
			if elm.Kind() == reflect.Ptr && !elm.IsNil() && elm.Elem().Kind() == reflect.Ptr {
				field = elm
			}
		}

		// if nested field is a pointer,
		// then extract it using .Elem()
		if field.Kind() == reflect.Ptr {
			field = field.Elem()
		}

		// fieldType.Tag.Get("json") return "", then it should be anonymous or embedded struct
		if (field.Kind() == reflect.Struct) && (fieldType.Anonymous || fieldType.Tag.Get("json") == "") {

			nesteds := hideFields(reflect.TypeOf(field.Interface()), reflect.ValueOf(field.Interface()), fs)
			for k, v := range nesteds {
				output[k] = v
			}
		} else {
			// get json Tag
			// eg: Name `json:"name"` => result will be `name`
			jsonField := fieldType.Tag.Get("json")
			if !fs[jsonField] {
				output[jsonField] = v.Field(i).Interface()
			}
		}

	}

	return output
}
