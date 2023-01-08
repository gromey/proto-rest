# RoundTripper

### The `roundtripper` package contains functions that are used as middleware on the client side.

## Getting Started

```go
package main

import (
	"net/http"

	"github.com/gromey/proto-rest/logger"
	"github.com/gromey/proto-rest/roundtripper"
)

func main() {
	rt := roundtripper.Sequencer(
		http.DefaultTransport,
		roundtripper.DumpHttp(logger.LevelTrace),
		roundtripper.Timer(logger.LevelInfo),
		roundtripper.PanicCatcher,
	)

	hClt := new(http.Client)
	hClt.Transport = rt
}
```