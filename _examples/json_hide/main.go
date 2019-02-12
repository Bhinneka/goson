package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/Bhinneka/goson"
)

// Bhinnekaners type
type Bhinnekaners struct {
	ID        string    `json:"id"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	BirthDate time.Time `json:"birthDate"`
	Created   time.Time `json:"created"`
	CreatorIP string    `json:"creatorIp"`
}

func main() {
	wury := &Bhinnekaners{
		ID:        "001",
		FirstName: "Wuriyanto",
		LastName:  "Musobar",
		BirthDate: time.Now(),
		Created:   time.Now(),
		CreatorIP: "10.0.1.1",
	}

	b1, err := json.Marshal(wury)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(string(b1))

	result, err := goson.HideFields(*wury, "created", "creatorIp")

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
