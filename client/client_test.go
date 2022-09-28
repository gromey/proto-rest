package client_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gromey/proto-rest/client"
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

type exampleStructSrv struct {
	Field int
}

type exampleStructClt struct {
	Field string
}

func makeTestSrv(t *testing.T, in any, out any) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		equal(t, r.URL.String(), "/path")

		if in != nil {
			input := &exampleStructSrv{}

			err := json.NewDecoder(r.Body).Decode(input)
			equal(t, nil, err)
			equal(t, in, input)
		}

		err := json.NewEncoder(w).Encode(out)
		equal(t, nil, err)
	}))
}

func TestProtoClient_Request(t *testing.T) {
	var tests = []struct {
		name   string
		method string
		input  any
		output any
		err    error
	}{
		{
			name:   "successful GET request",
			method: http.MethodGet,
			output: &exampleStructClt{Field: "example"},
		},
		{
			name:   "successful POST request",
			method: http.MethodPost,
			input:  &exampleStructSrv{Field: 1},
			output: &exampleStructClt{Field: "example"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			srv := makeTestSrv(t, test.input, test.output)
			defer srv.Close()

			clt := client.New(coder.NewCoder("application/json", json.Marshal, json.Unmarshal), srv.Client())

			resp, err := clt.Request(context.Background(), test.method, srv.URL+"/path", test.input, nil)
			equal(t, nil, err)

			defer func() { _ = resp.Body.Close() }()

			output := &exampleStructClt{}

			err = clt.Decode(resp.Body, output)
			equal(t, nil, err)
			equal(t, test.output, output)
		})
	}
}
