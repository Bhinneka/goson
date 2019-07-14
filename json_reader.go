package goson

import (
	"encoding/json"
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
	Uint() uint
	Int() int
	String() string
	Float() float64
	Bool() bool
	SetTo(target interface{}) error
}

// Read from json input
func Read(data []byte) Reader {
	var tmp interface{}
	json.Unmarshal(data, &tmp)

	return &reader{
		data: tmp,
	}
}

type reader struct {
	data interface{}
}

func (r *reader) GetField(keys ...string) Field {
	return nil
}

type field struct {
	value interface{}
}

func (f *field) Uint() uint {
	return 0
}

func (f *field) Int() int {
	return 0
}

func (f *field) Float() float64 {
	return 0.0
}

func (f *field) String() string {
	return ""
}

func (f *field) SetTo(target interface{}) error {
	return nil
}
