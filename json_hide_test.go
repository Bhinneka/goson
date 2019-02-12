package goson

import (
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

	wury := bhinnekaners{
		ID:        "001",
		FirstName: "Wuriyanto",
		LastName:  "Musobar",
		BirthDate: time.Now(),
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
	})
}
