package goson

import (
	"fmt"
	"reflect"
	"strings"
)

func makeZeroField(obj reflect.Value) {
	for i := 0; i < obj.NumField(); i++ {
		field := obj.Field(i)

		// if field is struct or types (nested struct or slice), process with recursive
		switch field.Kind() {
		case reflect.Ptr:
			processPointer(field)
		case reflect.Struct:
			makeZeroField(field)
		case reflect.Slice:
			fetchSliceType(field)
		}

		jsonTag := obj.Type().Field(i).Tag.Get("json")
		jsonTags := strings.Split(jsonTag, ",")
		if len(jsonTags) > 1 && jsonTags[1] == "omitempty" {
			field.Set(reflect.Zero(reflect.TypeOf(field.Interface())))
		}
	}
}

func processPointer(ptr reflect.Value) {
	val := ptr.Interface()
	if ptr.IsNil() {
		ptr = reflect.ValueOf(&val).Elem()
		ptr = reflect.New(ptr.Elem().Type().Elem()) // create from new domain model type of field
	}
	makeZeroField(ptr.Elem())
}

func fetchSliceType(slice reflect.Value) {
	for i := 0; i < slice.Len(); i++ {
		obj := slice.Index(i)
		if obj.Kind() == reflect.Ptr {
			processPointer(obj)
		} else {
			makeZeroField(obj)
		}
	}
}

func fetchMapType(mapper reflect.Value) {
	for _, idx := range mapper.MapKeys() {
		obj := mapper.MapIndex(idx)
		switch obj.Kind() {
		case reflect.Slice:
			fetchSliceType(obj)
		case reflect.Ptr:
			processPointer(obj)
		default:
			makeZeroField(obj)
		}
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
		processPointer(refValue)
	default:
		err = fmt.Errorf("invalid type %v: accept pointer, slice, and map", refValue.Kind())
	}

	return
}
