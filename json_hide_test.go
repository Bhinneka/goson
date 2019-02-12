package goson

import (
	"encoding/json"
	"testing"
	"time"
)

// json_hide test

type bhinnekaners struct {
	ID        string    `json:"id"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	BirthDate time.Time `json:"birthDate"`
	Created   time.Time `json:"created"`
	CreatorIP string    `json:"creatorIp"`
}

func TestJsonHide(t *testing.T) {

	birthDateLayout := "02/01/2006"
	birthDate, _ := time.Parse(birthDateLayout, "01/01/2019")

	wury := bhinnekaners{
		ID:        "001",
		FirstName: "Wuriyanto",
		LastName:  "Musobar",
		BirthDate: birthDate,
		Created:   time.Now(),
		CreatorIP: "10.0.1.1",
	}

	t.Run("should return choosen json field", func(t *testing.T) {
		result, err := HideFields(wury, "created", "creatorIp")

		if err != nil {
			t.Error("error hiding json fields")
		}

		if result == nil {
			t.Error("error hiding json fields")
		}

		if result != nil {
			t.Log(result)
		}

		expectedOutput := `{"birthDate":"2019-01-01T00:00:00Z","firstName":"Wuriyanto","id":"001","lastName":"Musobar"}`

		b, _ := json.Marshal(result)

		eq := string(b) == expectedOutput

		t.Log(string(b))

		if !eq {
			t.Error("result is not equal expectedOutput")
		}

	})

	t.Run("should return error if input is not struct", func(t *testing.T) {
		result, err := HideFields("", "created", "creatorIp")

		if err == nil {
			t.Error("should error, input is not struct")
		}

		if result != nil {
			t.Error("should error, input is not struct")
		}
	})
}
