package goson

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

var (
	// check basic data type (integer, float, string, boolean)
	intCheck = map[reflect.Kind]bool{
		reflect.Int8: true, reflect.Int16: true, reflect.Int32: true, reflect.Int: true, reflect.Int64: true,
	}
	floatCheck = map[reflect.Kind]bool{
		reflect.Float32: true, reflect.Float64: true,
	}
	stringCheck = map[reflect.Kind]bool{
		reflect.String: true,
	}
	boolCheck = map[reflect.Kind]bool{
		reflect.Bool: true,
	}
)

// Unmarshal json string, avoid error if incompatible data type
func Unmarshal(data []byte, target interface{}) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic: %v", r)
		}
	}()

	refValue := reflect.ValueOf(target)
	if ok := whiteListType[refValue.Kind().String()]; !ok {
		return fmt.Errorf("invalid target type %v: accept pointer, slice, and map", refValue.Kind())
	}

	var tmp map[string]interface{}
	err = json.Unmarshal(data, &tmp)
	if err != nil {
		return err
	}

	fetchDataType(refValue, tmp)
	return
}

func fetchDataType(obj reflect.Value, source map[string]interface{}) {
	switch obj.Kind() {
	case reflect.Struct:
		scan(obj, source)

	case reflect.Ptr:
		fetchDataType(obj.Elem(), source)
	}
}

// scan only if data type is struct
func scan(obj reflect.Value, data map[string]interface{}) {
	objType := obj.Type()
	for i := 0; i < obj.NumField(); i++ {
		field := obj.Field(i)
		if !field.CanSet() { // break if field cannot set a value (usually an unexported field in struct), to avoid a panic
			return
		}

		fetchDataType(field, data)

		jsonTag := objType.Field(i).Tag.Get("json")
		if jsonTags := strings.Split(jsonTag, ","); len(jsonTags) > 0 {
			jsonTag = jsonTags[0]
		}
		if jsonTag == "" {
			jsonTag = objType.Field(i).Name
		}

		setValue(field, data[jsonTag])
	}
}

func setValue(targetField reflect.Value, data interface{}) {
	if !targetField.IsValid() {
		return
	}

	targetKind := targetField.Kind()     // targetKind is datatype from target
	valueSource := reflect.ValueOf(data) // valueSource is datatype from json source

	// check target is pointer or not, and value from json data source
	if targetKind == reflect.Ptr {
		if data != nil && targetField.IsNil() {
			rv := reflect.ValueOf(targetField.Interface())
			val := reflect.New(rv.Type().Elem()).Interface()
			targetField.Set(reflect.ValueOf(val))
		}
		targetField = targetField.Elem() // take the element if target is pointer, to set a value in target
		targetKind = targetField.Kind()
	}

	// switch datatype from json source
	switch valueSource.Kind() {
	case reflect.String: // field from json source is string
		str := valueSource.Interface().(string)
		switch {
		case stringCheck[targetKind]:
			targetField.Set(valueSource)
		case intCheck[targetKind]:
			if val, err := strconv.Atoi(str); err == nil {
				targetField.Set(reflect.ValueOf(int(val)))
			}
		case floatCheck[targetKind]:
			if val, err := strconv.ParseFloat(str, -1); err == nil {
				value := reflect.ValueOf(val)
				if targetKind == reflect.Float32 {
					value = reflect.ValueOf(float32(val))
				}
				targetField.Set(value)
			}
		case boolCheck[targetKind]:
			if val, err := strconv.ParseBool(str); err == nil {
				targetField.Set(reflect.ValueOf(val))
			}
		}

	case reflect.Float64: // field from json source is float, and integer (because integer in json source will be made to Float64 when Unmarshal)
		fl := valueSource.Interface().(float64)
		switch {
		case floatCheck[targetKind]:
			if targetKind == reflect.Float32 {
				valueSource = reflect.ValueOf(float32(fl))
			}
			targetField.Set(valueSource)
		case intCheck[targetKind]:
			targetField.Set(reflect.ValueOf(int(fl)))
		case stringCheck[targetKind], boolCheck[targetKind]:
			targetField.Set(reflect.ValueOf(strconv.FormatFloat(fl, 'f', -1, 64)))
		case boolCheck[targetKind]:
			if v, err := strconv.ParseBool(strconv.FormatFloat(fl, 'f', -1, 64)); err == nil {
				targetField.Set(reflect.ValueOf(v))
			}
		}

	case reflect.Bool: // field from json source is boolean
		switch {
		case boolCheck[targetKind]:
			targetField.Set(valueSource)
		case stringCheck[targetKind]:
			bl := valueSource.Interface().(bool)
			targetField.Set(reflect.ValueOf(strconv.FormatBool(bl)))
		}

	case reflect.Map: // representation from struct, process with recursive again
		subData := valueSource.Interface().(map[string]interface{})
		fetchDataType(targetField, subData)

	case reflect.Slice: // representation from array, process with recursive again
		if targetKind == reflect.Slice {
			data := valueSource.Interface().([]interface{})
			tmpSlice := reflect.MakeSlice(reflect.SliceOf(reflect.TypeOf(targetField.Interface()).Elem()), len(data), len(data))
			for i := 0; i < tmpSlice.Len(); i++ {
				setValue(tmpSlice.Index(i), data[i])
			}
			targetField.Set(tmpSlice)
		}
	}
}
