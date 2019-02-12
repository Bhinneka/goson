package main

import (
	"encoding/json"
	"fmt"

	"github.com/Bhinneka/goson"
)

// Child type
type Child struct {
	ChildA      string `json:"childA"`
	ChildB      string `json:"childB,omitempty"`
	ChildStruct struct {
		A int `json:"A,omitempty"`
		B int `json:"B"`
	} `json:"childStruct"`
}

// Bhinnekaners type
type Bhinnekaners struct {
	Name       string `json:"name"`
	Additional string `json:"additional,omitempty"`
	AddStruct  struct {
		Field1 string `json:"field1,omitempty"`
		Field2 string `json:"field2"`
	} `json:"addStruct"`
	AddChildPtr *Child `json:"child,omitempty"`
}

func main() {
	objExample := Bhinnekaners{
		Name:       "Agung DP",
		Additional: "test lagi",
		AddStruct: struct {
			Field1 string `json:"field1,omitempty"`
			Field2 string `json:"field2"`
		}{
			Field1: "Field 1",
			Field2: "Field 2",
		},
		AddChildPtr: &Child{
			ChildA: "Child A",
			ChildB: "Child B",
			ChildStruct: struct {
				A int `json:"A,omitempty"`
				B int `json:"B"`
			}{
				A: 453,
				B: 567,
			},
		},
	}

	jsonBefore, err := json.Marshal(objExample)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(jsonBefore))
	// will print -> {"name":"Agung DP","additional":"test lagi","addStruct":{"field1":"Field 1","field2":"Field 2"},"child":{"childA":"Child A","childB":"Child B","childStruct":{"A":453,"B":567}}}

	err = goson.MakeZeroOmitempty(&objExample)
	if err != nil {
		panic(err)
	}

	jsonAfter, err := json.Marshal(objExample)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(jsonAfter))
	// will print -> {"name":"Agung DP","addStruct":{"field2":"Field 2"},"child":{"childA":"Child A","childStruct":{"B":567}}}
}
