### Go JSON utility library

[![GoDoc](https://godoc.org/github.com/Bhinneka/goson?status.svg)](https://godoc.org/github.com/Bhinneka/goson)
[![Go Report Card](https://goreportcard.com/badge/github.com/Bhinneka/goson)](https://goreportcard.com/report/github.com/Bhinneka/goson)

#### Install
```shell
$ go get github.com/Bhinneka/goson
```

#### Todo
- Add more utility

### Hide specific fields
```go
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

    // hide specific field, eg: fields created and creatorIp
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
}

```

<b>before hide</b>
```json
{
	"id": "001",
	"firstName": "Wuriyanto",
	"lastName": "Musobar",
	"birthDate": "2019-02-12T11:07:20.68801+07:00",
	"created": "2019-02-12T11:07:20.68801+07:00",
	"creatorIp": "10.0.1.1"
}
```
<b>after hide</b>
```json
{
	"id": "001",
	"lastName": "Musobar",
	"firstName": "Wuriyanto",
	"birthDate": "2019-02-12T11:12:04.5998+07:00"
}
```

### Make Zero Omitempty

```go
// Child type
type Child struct {
	ChildA      string `json:"childA"`
	ChildB      string `json:"childB,omitempty"`
	ChildStruct struct {
		A int `json:"A,omitempty"`
		B int `json:"B"`
	} `json:"childStruct"`
}
// arr type for slice example
type arr struct {
	Arr1 string `json:"arr1"`
	Arr2 string `json:"arr2,omitempty"`
}

// Bhinnekaners type
type Bhinnekaners struct {
	Name       string `json:"name"`
	Additional string `json:"additional,omitempty"`
	AddStruct  struct {
		Field1 string `json:"field1,omitempty"`
		Field2 string `json:"field2"`
	} `json:"addStruct"`
	AddChildPtr *Child `json:"child"`
	Slice       []arr  `json:"slice"`
}

example := Bhinnekaners{
	Name:       "Agung DP",
	Additional: "testing",
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
	Slice: []arr{
		arr{
			Arr1: "Test field slice 1",
			Arr2: "Test field slice 2",
		},
	},
}
```

<b>json marshal before make zero field omitempty in example variable:</b>
```json
{
  "name": "Agung DP",
  "additional": "testing",
  "addStruct": {
    "field1": "Field 1",
    "field2": "Field 2"
  },
  "child": {
    "childA": "Child A",
    "childB": "Child B",
    "childStruct": {
      "A": 453,
      "B": 567
    }
  },
  "slice": [
    {
      "arr1": "Test field slice 1",
      "arr2": "Test field slice 2"
    }
  ]
}
```

<b>json marshal after make zero field omitempty in example variable:</b>
```json
{
  "name": "Agung DP",
  "addStruct": {
    "field2": "Field 2"
  },
  "child": {
    "childA": "Child A",
    "childStruct": {
      "B": 567
    }
  },
  "slice": [
    {
      "arr1": "Test field slice 1"
    }
  ]
}
```


### Advance JSON Unmarshal

```go
package main

import (
	"fmt"
	"github.com/Bhinneka/goson"
)

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

var data = []byte(`
{
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
}
`)

func main() {
	var target Model
	goson.Unmarshal(data, &target)
	fmt.Printf("%+v\n", target)
}
```

<b>Value from target:</b>

```json
{
     "id": 1,
     "name": "agungdp",
     "mustFloat": 2.23423,
     "mustInt": 2,
     "isExist": true,
     "obj": {
          "n": 2,
          "testing": {
               "ss": 23840923849284
          }
     },
     "slice": [
          {
               "fieldA": 100,
               "fieldB": "3000",
               "exist": "true"
          },
          {
               "fieldA": 50000,
               "fieldB": "123000",
               "exist": ""
          }
     ]
}
```