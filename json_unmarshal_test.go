package goson

import (
	"fmt"
	"testing"
)

func TestUnmarshal(t *testing.T) {
	type Slice struct {
		FieldA uint16  `json:"fieldA"`
		FieldB string  `json:"fieldB"`
		Exist  string  `json:"exist"`
		Test   float32 `json:"test"`
	}
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
		skip    string // cannot set value to this field because unexported
	}

	t.Run("Testcase #1: Testing Unmarshal with root is JSON Object", func(t *testing.T) {
		data := []byte(`{
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
		var target Model
		if err := Unmarshal(data, &target); err != nil {
			t.Errorf("Unmarshal() error = %v", err)
		}
		if target.MustFloat == nil {
			t.Errorf("field should not nil")
		}
		if target.Uint != 11 {
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
		if target.Bools[2] != false {
			t.Errorf("value not equal")
		}
		if target.NoTag != "19283091832" {
			t.Errorf("value not equal")
		}

		fmt.Printf("%+v\n\n", target)
	})

	t.Run("Testcase #2: Testing Unmarshal with root is JSON Array", func(t *testing.T) {
		data := []byte(`[
			{
				 "fieldA": 100,
				 "fieldB": "3000",
				 "exist": "true",
				 "test": 3.14
			},
			{
				 "fieldA": 50000,
				 "fieldB": "123000",
				 "exist": "3000",
				 "test": 1323.123
			}
	   ]`)
		var target []Slice
		if err := Unmarshal(data, &target); err != nil {
			t.Errorf("Unmarshal() error = %v", err)
		}
		if target == nil {
			t.Errorf("field should not nil")
		}
		if len(target) != 2 {
			t.Errorf("not equal")
		}

		fmt.Printf("%+v\n\n", target)
	})

	t.Run("Testcase #3: Testing error unmarshal (target is not pointer)", func(t *testing.T) {
		data := []byte(`{}`)
		var target Model
		err := Unmarshal(data, target)
		if err == nil {
			t.Errorf("Unmarshal() must error, got nil")
		}
	})

	t.Run("Testcase #4: Testing error unmarshal (invalid json format)", func(t *testing.T) {
		data := []byte(`{1: satu}`)
		var target Model
		err := Unmarshal(data, &target)
		if err == nil {
			t.Errorf("Unmarshal() must error, got nil")
		}
	})
}
