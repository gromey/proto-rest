package coder_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"reflect"
	"testing"

	"github.com/gromey/proto-rest/coder"
	"github.com/gromey/proto-rest/logger"
)

func init() {
	logger.SetLogger(logger.New(nil))
}

func equal(t *testing.T, exp, got any) {
	if !reflect.DeepEqual(exp, got) {
		t.Fatalf("Not equal:\nexp: %v\ngot: %v", exp, got)
	}
}

type exampleStruct struct {
	Field string `json:"field"`
}

func TestEncoder_Encode(t *testing.T) {
	encoder := coder.NewEncoder(json.Marshal)
	var tests = []struct {
		name   string
		input  *exampleStruct
		output []byte
		err    error
	}{
		{
			name:   "successful encode",
			input:  &exampleStruct{Field: "example"},
			output: []byte("{\"field\":\"example\"}"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			err := encoder.Encode(w, test.input)
			if test.err != nil {
				equal(t, test.err.Error(), err.Error())
			} else {
				equal(t, nil, err)
				equal(t, test.output, w.Bytes())
			}
		})
	}
}

func TestDecoder_Decode(t *testing.T) {
	decoder := coder.NewDecoder(json.Unmarshal)
	var tests = []struct {
		name   string
		input  []byte
		output *exampleStruct
		err    error
	}{
		{
			name:   "successful decode",
			input:  []byte("{\"field\":\"example\"}"),
			output: &exampleStruct{Field: "example"},
		},
		{
			name:  "unexpected end of JSON input",
			input: []byte("{\"field\":\"example\""),
			err:   errors.New("unexpected end of JSON input"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r := &bytes.Buffer{}
			r.Write(test.input)

			v := new(exampleStruct)

			err := decoder.Decode(r, v)
			if test.err != nil {
				equal(t, test.err.Error(), err.Error())
			} else {
				equal(t, nil, err)
				equal(t, test.output, v)
			}
		})
	}
}
