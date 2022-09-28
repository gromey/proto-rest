# Coder

### The `coder` package implements three interfaces with debug [logging](https://github.com/gromey/proto-rest/blob/main/logger/README.md):

- *Encoder* encodes and writes values to an output stream.
- *Decoder* reads and decodes values from an input stream.
- *Coder* is a pair of Encoder and Decoder.

## Getting Started

```go
package main

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/gromey/proto-rest/coder"
)

func main() {
	coderJSON := coder.NewCoder("application/json", json.Marshal, json.Unmarshal)

	var buf bytes.Buffer

	in := &struct {
		A string `json:"a"`
	}{A: "AAA"}

	if err := coderJSON.Encode(&buf, in); err != nil {
		panic(err)
	}

	fmt.Printf("encoded: %s\n", in)
	// [DEBUG] Func: Encode() Encoder input data: &struct { A string "json:\"a\"" }{A:"AAA"}
	// [DEBUG] Func: Encode() Encoder output data: {"a":"AAA"}
	// encoded: {"a":"AAA"}

	out := &struct {
		A string `json:"a"`
	}{}

	if err := coderJSON.Decode(&buf, out); err != nil {
		panic(err)
	}

	fmt.Printf("decoded: %+v\n", out)
	// [DEBUG] Func: Decode() Decoder input data: {"a":"AAA"}
	// [DEBUG] Func: Decode() Decoder output data: &struct { A string "json:\"a\"" }{A:"AAA"}
	// decoded: &{A:AAA}
}
```