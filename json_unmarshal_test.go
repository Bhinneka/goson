package goson

import (
	"testing"
)

func TestUnmarshal(t *testing.T) {
	type Slice struct {
		FieldA int     `json:"fieldA"`
		FieldB string  `json:"fieldB"`
		Exist  string  `json:"exist"`
		Test   float32 `json:"test"`
	}
	type Model struct {
		ID        int      `json:"id" default:"1"`
		Name      string   `json:"name"`
		MustFloat *float64 `json:"mustFloat"`
		MustInt   int      `json:"mustInt"`
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
	}

	t.Run("Testing unmarshal", func(t *testing.T) {
		data := []byte(`{
				"id": "01",
				"name": "agungdp",
				"mustFloat": "2.23423",
				"mustInt": 2.23423,
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
					"fieldB": "3000",
					"exist": true
				  },
				  {
					"fieldA": 50000,
					"fieldB": "123000",
					"test": 1323.123
				  }
				],
				"str": ["a", true],
				"ints": ["2", 3]
			  }`)
		var target Model
		if err := Unmarshal(data, &target); err != nil {
			t.Errorf("Unmarshal() error = %v", err)
		}
		if target.MustFloat == nil {
			t.Errorf("field should not nil")
		}
		if target.Obj.Testing.Ss != 23840923849284 {
			t.Errorf("value not equal")
		}
		if target.Slice[0].Exist != "true" {
			t.Errorf("value not equal")
		}
		if *target.Strings[1] != "true" {
			t.Errorf("value not equal")
		}
		if target.Ints[0] != 2 {
			t.Errorf("value not equal")
		}

		// fmt.Println("==================================================================")
		// fmt.Printf("\n%+v\n\n", target)
	})
}
