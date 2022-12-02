# Client

### The `client` package implements the core functionality of a [coder](https://github.com/gromey/proto-rest/blob/main/coder/README.md)-based client with trace [logging](https://github.com/gromey/proto-rest/blob/main/logger/README.md).

## Getting Started

```go
package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/gromey/proto-rest/client"
	"github.com/gromey/proto-rest/coder"
	"github.com/gromey/proto-rest/logger"
)

func main() {
	coderJSON := coder.NewCoder("application/json", json.Marshal, json.Unmarshal)

	clientJSON := client.New(coderJSON, http.DefaultClient)

	cUrl, err := url.Parse("http://localhost:8080/example/")
	if err != nil {
		panic(err)
	}

	params := make(url.Values)
	params.Add("id", "1")
	cUrl.RawQuery = params.Encode()

	// To add additional data to the request, use the optional function f(*http.Request)
	f := func(r *http.Request) {
		r.Header.Set("Accept", "application/json")
	}

	resp, err := clientJSON.Request(context.TODO(), http.MethodGet, cUrl.String(), nil, f)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	res := &struct{}{}

	err = clientJSON.Decode(resp.Body, res)
	if err != nil {
		panic(err)
	}
}
```

```go
package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gromey/proto-rest/client"
	"github.com/gromey/proto-rest/coder"
	"github.com/gromey/proto-rest/logger"
)

func main() {
	coderJSON := coder.NewCoder("application/json", json.Marshal, json.Unmarshal)

	clientJSON := client.New(coderJSON, http.DefaultClient)

	cUrl := "http://localhost:8080/v1/example/"

	req := &struct {
		ID int `json:"id"`
	}{ID: 1}

	resp, err := clientJSON.Request(context.TODO(), http.MethodPost, cUrl, req, nil)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	res := &struct {
		// some fields
	}{}

	err = clientJSON.Decode(resp.Body, res)
	if err != nil {
		panic(err)
	}
}
```