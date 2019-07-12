package main

import (
	"fmt"

	"github.com/Bhinneka/goson"
)

// Slice struct
type Slice struct {
	FieldA uint16  `json:"fieldA"`
	FieldB string  `json:"fieldB"`
	Exist  string  `json:"exist"`
	Test   float32 `json:"test"`
}

// Model struct
type Model struct {
	ID        int      `json:"id"`
	Name      string   `json:"name"`
	MustFloat *float64 `json:"mustFloat"`
	MustInt   int      `json:"mustInt"`
	Uint      uint     `json:"uint"`
	IsExist   *bool    `json:"isExist"`
	Obj       *struct {
		N       int `json:"n"`
		Testing struct {
			Ss int `json:"ss"`
		} `json:"testing"`
	} `json:"obj"`
	Slice   []Slice   `json:"slice"`
	Strings []*string `json:"str"`
	Ints    []int     `json:"ints"`
	Bools   []bool    `json:"bools"`
	NoTag   string
}

var data = []byte(`{
	"id": "01",
	"name": "agungdp",
	"mustFloat": "2.23423",
	"mustInt": 2.23423,
	"uint": 11,
	"isExist": "true",
	"obj": {
	  "n": 2,
	  "testing": {
		"ss": "23840923849284"
	  }
	},
	"slice": [
	  {
		"fieldA": "100",
		"fieldB": 3000,
		"exist": true,
		"test": "3.14"
	  },
	  {
		"fieldA": 50000,
		"fieldB": "123000",
		"exist": 3000,
		"test": 1323.123
	  }
	],
	"str": ["a", true],
	"ints": ["2", 3],
	"bools": [1, "true", "0", true],
	"NoTag": 19283091832
}`)

func main() {
	var model Model
	err := goson.Unmarshal(data, &model)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", model)
}
