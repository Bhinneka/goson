package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/Bhinneka/goson"
)

// this example simulate hide fields from embedded/ nested and pointer of struct

// Bhinnekaners type
type Bhinnekaners struct {
	ID        string    `json:"id"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	BirthDate time.Time `json:"birthDate"`
	*Base1
}

// Base1 struct
type Base1 struct {
	Created   time.Time `json:"created"`
	CreatorIP string    `json:"creatorIp"`
	Base2
}

// Base2 struct
type Base2 struct {
	CreatorID string `json:"creatorId"`
	EditorID  string `json:"editorId"`
}

func main() {
	wury := &Bhinnekaners{
		ID:        "001",
		FirstName: "Wuriyanto",
		LastName:  "Musobar",
		BirthDate: time.Now(),
		Base1: &Base1{
			Created:   time.Now(),
			CreatorIP: "10.0.1.1",
			Base2: Base2{
				CreatorID: "1900",
				EditorID:  "1800",
			},
		},
	}

	b1, err := json.Marshal(wury)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(string(b1))

	result, err := goson.HideFields(*wury, "editorId")

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	b2, err := json.Marshal(result)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(string(b2))
	fmt.Println()

	json.NewEncoder(os.Stdout).Encode(result)
}
