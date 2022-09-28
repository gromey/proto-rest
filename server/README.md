# Server

### The `server` package implements the core functionality of a [coder](https://github.com/gromey/proto-rest/blob/main/coder/README.md)-based server.

## Getting Started

```go
package main

import (
	"encoding/json"
	"net/http"

	"github.com/gromey/proto-rest/coder"
	"github.com/gromey/proto-rest/server"
)

func main() {
	coderJSON := coder.NewCoder("application/json", json.Marshal, json.Unmarshal)

	serverJSON := server.New(coderJSON)

	handlerFunc := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			req := &struct {
				// some fields
			}{}

			if err := serverJSON.Decode(r.Body, req); err != nil {
				panic(err)
			}
		}

		res := &struct {
			ID int `json:"id"`
		}{ID: 1}

		serverJSON.WriteResponse(w, http.StatusOK, res)
	})

	http.Handle("/example/", handlerFunc)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
```

### For all responses with a body:

- The default is to use the `Content-Type` set in
  the [coder](https://github.com/gromey/proto-rest/blob/main/coder/README.md).
- If you don't set the `Content-Type` in the [coder](https://github.com/gromey/proto-rest/blob/main/coder/README.md), it
  will be set automatically by the [net/http](https://pkg.go.dev/net/http) package.
- If you need to set a different `Content-Type` you must set it before calling `WriteResponse`.

### For all responses without a body:

- `Content-Type` will not be set by default.
- If you need to set `Content-Type` you must set it before calling `WriteResponse`.