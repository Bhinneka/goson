package goson

import (
	"fmt"
	"reflect"
	"strings"
)

func makeZeroField(obj reflect.Value) {
	for i := 0; i < obj.NumField(); i++ {
		isFieldPtr := false // for sign field is pointer or not

		fieldValue := obj.Field(i)
		if fieldValue.Kind() == reflect.Ptr && !fieldValue.IsNil() {
			fieldValue = fieldValue.Elem()
			isFieldPtr = true
		}
		if fieldValue.Kind() == reflect.Struct {
			makeZeroField(fieldValue) // if field is struct or types (nested struct), process with recursive
		}

		jsonTag := obj.Type().Field(i).Tag.Get("json")
		jsonTags := strings.Split(jsonTag, ",")
		if !isFieldPtr && len(jsonTags) > 1 && jsonTags[1] == "omitempty" {
			fieldValue.Set(reflect.Zero(reflect.TypeOf(fieldValue.Interface())))
		}
	}
}

func fetchSliceType(slice reflect.Value) {
	for i := 0; i < slice.Len(); i++ {
		obj := slice.Index(i)
		if obj.Kind() == reflect.Ptr {
			obj = obj.Elem()
		}
		makeZeroField(obj)
	}
}

func fetchMapType(mapper reflect.Value) {
	for _, idx := range mapper.MapKeys() {
		obj := mapper.MapIndex(idx)
		if obj.Kind() == reflect.Ptr {
			obj = obj.Elem()
		} else if obj.Kind() == reflect.Slice {
			fetchSliceType(obj)
			continue
		}
		makeZeroField(obj)
	}
}

// MakeZeroOmitempty tools for make Zero value in field contains `json: omitempty` tag
func MakeZeroOmitempty(obj interface{}) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()

	refValue := reflect.ValueOf(obj)

	switch refValue.Kind() {
	case reflect.Slice:
		fetchSliceType(refValue)
	case reflect.Map:
		fetchMapType(refValue)
	case reflect.Ptr:
		makeZeroField(refValue.Elem())
	default:
		err = fmt.Errorf("invalid type %v: accept pointer, slice, and map", refValue.Kind())
	}

	return
}
