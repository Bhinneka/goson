package goson

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

var (
	arrayCheck = map[byte]byte{
		'[': ']',
	}
)

// Reader json string
type Reader interface {
	GetField(keys ...string) Field
}

// Field from json reader
type Field interface {
	Uint() (uint, error)
	Int() (int, error)
	String() (string, error)
	Float() (float64, error)
	Bool() (bool, error)
	SetTo(target interface{}) error
}

// Read from json input
func Read(data []byte) (Reader, error) {
	var tmp interface{}
	err := json.Unmarshal(data, &tmp)
	if err != nil {
		return nil, err
	}

	return &reader{
		data: tmp,
	}, nil
}

type reader struct {
	data interface{}
}

func (r *reader) GetField(keys ...string) Field {
	var data interface{} = r.data
	field := new(field)
	func() {
		var traces []string
		defer func() {
			if r := recover(); r != nil {
				field.err = fmt.Errorf("key '%s' not found or out of range from JSON data", strings.Join(traces, " => "))
			}
		}()
		for _, k := range keys {
			traces = append(traces, k)
			if arrayCheck[k[0]] == k[len(k)-1] {
				if idx, err := strconv.Atoi(k[1 : len(k)-1]); err == nil {
					data = data.([]interface{})[idx]
				}
				continue
			}
			tmp, ok := data.(map[string]interface{})[k]
			if !ok && tmp == nil {
				panic("invalid key")
			}
			data = tmp
		}
	}()

	field.value = data
	return field
}

type field struct {
	err   error
	value interface{}
}

func (f *field) set(v reflect.Value) error {
	if f.err != nil {
		return f.err
	}
	return setValue(v, f.value)
}

func (f *field) Uint() (u uint, err error) {
	err = f.set(reflect.ValueOf(&u))
	return
}

func (f *field) Int() (i int, err error) {
	err = f.set(reflect.ValueOf(&i))
	return
}

func (f *field) Float() (fl float64, err error) {
	err = f.set(reflect.ValueOf(&fl))
	return
}

func (f *field) String() (s string, err error) {
	err = f.set(reflect.ValueOf(&s))
	return
}

func (f *field) Bool() (b bool, err error) {
	err = f.set(reflect.ValueOf(&b))
	return
}

func (f *field) SetTo(target interface{}) (err error) {
	err = f.set(reflect.ValueOf(target))
	return
}
