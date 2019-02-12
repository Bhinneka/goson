### Go JSON utility library

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