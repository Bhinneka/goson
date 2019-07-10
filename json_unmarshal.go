package goson

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

var (
	// check basic type data (integer, float, string, boolean)
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
		if obj.IsNil() {
			val := obj.Interface()
			rv := reflect.ValueOf(val)
			val = reflect.New(rv.Type().Elem()).Interface()
			obj.Set(reflect.ValueOf(val))
		}
		fetchDataType(obj.Elem(), source)
	}
}

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
	if data == nil || !targetField.IsValid() {
		return
	}

	targetKind := targetField.Kind()     // targetKind is datatype from target
	valueSource := reflect.ValueOf(data) // valueSource is datatype from json source

	// switch datatype from json source
	switch valueSource.Kind() {
	case reflect.String:
		str := valueSource.Interface().(string)
		if stringCheck[targetKind] {
			targetField.Set(valueSource)
		} else if intCheck[targetKind] {
			if val, err := strconv.Atoi(str); err == nil {
				targetField.Set(reflect.ValueOf(int(val)))
			}
		} else if floatCheck[targetKind] {
			if val, err := strconv.ParseFloat(str, -1); err == nil {
				targetField.Set(reflect.ValueOf(val))
			}
		} else if boolCheck[targetKind] {
			if val, err := strconv.ParseBool(str); err == nil {
				targetField.Set(reflect.ValueOf(val))
			}
		}

	case reflect.Float64:
		fl := valueSource.Interface().(float64)
		if floatCheck[targetKind] {
			targetField.Set(valueSource)
		} else if intCheck[targetKind] {
			targetField.Set(reflect.ValueOf(int(fl)))
		} else if stringCheck[targetKind] {
			targetField.Set(reflect.ValueOf(strconv.FormatFloat(fl, 'f', -1, 64)))
		}

	case reflect.Bool:
		if boolCheck[targetKind] {
			targetField.Set(valueSource)
		} else if stringCheck[targetKind] {
			bl := valueSource.Interface().(bool)
			targetField.Set(reflect.ValueOf(strconv.FormatBool(bl)))
		}

	case reflect.Map: // representation from struct
		subData := valueSource.Interface().(map[string]interface{})
		fetchDataType(targetField, subData)

	case reflect.Slice:
		if targetKind == reflect.Slice {
			data := valueSource.Interface().([]interface{})
			tmpSlice := reflect.MakeSlice(reflect.SliceOf(reflect.TypeOf(targetField.Interface()).Elem()), len(data), len(data))
			for i := 0; i < tmpSlice.Len(); i++ {
				subData := data[i].(map[string]interface{})
				obj := tmpSlice.Index(i)
				fetchDataType(obj, subData)
			}
			targetField.Set(tmpSlice)
		}
	}
}
