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
	uintCheck = map[reflect.Kind]bool{
		reflect.Uint: true, reflect.Uint8: true, reflect.Uint16: true, reflect.Uint32: true, reflect.Uint64: true,
	}
	intCheck = map[reflect.Kind]bool{
		reflect.Int: true, reflect.Int8: true, reflect.Int16: true, reflect.Int32: true, reflect.Int64: true,
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

	targetValue := reflect.ValueOf(target)
	if targetValue.Kind() != reflect.Ptr {
		return fmt.Errorf("invalid target type %v: must pass address from target", targetValue.Kind())
	}

	var dataSource interface{}
	err = json.Unmarshal(data, &dataSource)
	if err != nil {
		return err
	}

	scanTarget(targetValue, dataSource)
	return
}

// scanTarget from destination data for target type is non basic data type (int, float, string, bool)
func scanTarget(target reflect.Value, source interface{}) {
	switch target.Kind() {
	case reflect.Struct:
		data, _ := source.(map[string]interface{}) // if target is struct, source type must be map[string]interface{}
		scanStructField(target, data)

	case reflect.Slice:
		sourceVal := reflect.ValueOf(source)
		if sourceVal.Kind() == reflect.Slice {
			tmpSlice := reflect.MakeSlice(reflect.SliceOf(reflect.TypeOf(target.Interface()).Elem()), sourceVal.Len(), sourceVal.Len())
			for i := 0; i < tmpSlice.Len(); i++ {
				setValue(tmpSlice.Index(i), sourceVal.Index(i).Interface())
			}
			target.Set(tmpSlice)
		}

	case reflect.Ptr:
		scanTarget(target.Elem(), source)

	}
}

// scan only if data type is struct
func scanStructField(obj reflect.Value, data map[string]interface{}) {
	objType := obj.Type()
	for i := 0; i < obj.NumField(); i++ {
		field := obj.Field(i)
		if !field.CanSet() { // break if field cannot set a value (usually an unexported field in struct), to avoid a panic
			return
		}

		jsonTag := objType.Field(i).Tag.Get("json")
		if jsonTags := strings.Split(jsonTag, ","); len(jsonTags) > 0 { // if json tag contains "omitempty"
			jsonTag = jsonTags[0]
		}
		if jsonTag == "" {
			jsonTag = objType.Field(i).Name
		}

		source := data[jsonTag]
		scanTarget(field, source)

		setValue(field, source)
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
		if data != nil && targetField.IsNil() { // allocate new value to pointer target
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
			targetField.SetString(str)
		case intCheck[targetKind]:
			if val, err := strconv.Atoi(str); err == nil {
				targetField.SetInt(int64(val))
			}
		case uintCheck[targetKind]:
			if val, err := strconv.Atoi(str); err == nil {
				targetField.SetUint(uint64(val))
			}
		case floatCheck[targetKind]:
			if val, err := strconv.ParseFloat(str, -1); err == nil {
				targetField.SetFloat(val)
			}
		case boolCheck[targetKind]:
			if val, err := strconv.ParseBool(str); err == nil {
				targetField.SetBool(val)
			}
		}

	case reflect.Float64: // field from json source is float, and integer (because any integer in json source will be made to Float64 when Unmarshal)
		fl := valueSource.Interface().(float64)
		switch {
		case floatCheck[targetKind]:
			targetField.SetFloat(fl)
		case intCheck[targetKind]:
			targetField.SetInt(int64(fl))
		case uintCheck[targetKind]:
			targetField.SetUint(uint64(fl))
		case stringCheck[targetKind]:
			targetField.SetString(strconv.FormatFloat(fl, 'f', -1, 64))
		case boolCheck[targetKind]:
			if v, err := strconv.ParseBool(strconv.FormatFloat(fl, 'f', -1, 64)); err == nil {
				targetField.SetBool(v)
			}
		}

	case reflect.Bool: // field from json source is boolean
		bl := valueSource.Interface().(bool)
		switch {
		case boolCheck[targetKind]:
			targetField.SetBool(bl)
		case stringCheck[targetKind]:
			targetField.SetString(strconv.FormatBool(bl))
		}

	case reflect.Map, reflect.Slice: // representation from subtree in json source, process with recursive again
		scanTarget(targetField, data)

	}
}
