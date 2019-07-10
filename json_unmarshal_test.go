package goson

import (
	"fmt"
	"testing"
)

func TestUnmarshal(t *testing.T) {
	type Slice struct {
		FieldA int    `json:"fieldA"`
		FieldB string `json:"fieldB"`
		Exist  string `json:"exist"`
	}
	type Model struct {
		ID        int     `json:"id" default:"1"`
		Name      string  `json:"name"`
		MustFloat float64 `json:"mustFloat"`
		MustInt   int     `json:"mustInt"`
		IsExist   bool    `json:"isExist"`
		Obj       *struct {
			N       int `json:"n"`
			Testing struct {
				Ss int `json:"ss"`
			} `json:"testing"`
		} `json:"obj"`
		Slice []Slice `json:"slice"`
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
					"fieldB": "123000"
				  }
				]
			  }`)
		var target Model
		if err := Unmarshal(data, &target); err != nil {
			t.Errorf("Unmarshal() error = %v", err)
		}
		// fmt.Println("==================================================================")
		fmt.Printf("\n%+v\n\n", target)
	})
}
