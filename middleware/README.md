# Middleware

### The `middleware` package contains functions that are used as middleware on the server side.

## Getting Started

```go
package main

import (
	"fmt"
	"net/http"

	"github.com/gromey/proto-rest/logger"
	"github.com/gromey/proto-rest/middleware"
)

func main() {
	handlerFunc := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := fmt.Fprintln(w, "Hello World!"); err != nil {
			panic(err)
		}
	})

	http.Handle("/example/", handlerFunc)

	h := middleware.Sequencer(
		http.DefaultServeMux,
		middleware.DumpHttp(logger.LevelTrace),
		middleware.Timer(logger.LevelInfo),
		middleware.PanicCatcher,
	)

	if err := http.ListenAndServe(":8080", h); err != nil {
		panic(err)
	}
}
```