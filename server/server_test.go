package server_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gromey/proto-rest/coder"
	"github.com/gromey/proto-rest/logger"
	"github.com/gromey/proto-rest/server"
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

func makeRequest(t *testing.T, method string, srv *httptest.Server, in any, out any) {
	var reader io.Reader
	if in != nil {
		buf := new(bytes.Buffer)

		err := json.NewEncoder(buf).Encode(in)
		equal(t, nil, err)

		reader = buf
	}

	req, err := http.NewRequest(method, srv.URL+"/path", reader)
	equal(t, nil, err)

	resp, err := srv.Client().Do(req)
	equal(t, nil, err)

	defer func() { _ = resp.Body.Close() }()

	output := &exampleStructClt{}

	err = json.NewDecoder(resp.Body).Decode(output)
	equal(t, nil, err)
	equal(t, out, output)
}

func TestProtoServer_WriteResponse(t *testing.T) {
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
			srvJSON := server.New(coder.NewCoder("application/json", json.Marshal, json.Unmarshal))

			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				equal(t, r.URL.String(), "/path")

				if test.input != nil {
					input := &exampleStructSrv{}

					err := srvJSON.Decode(r.Body, input)
					equal(t, nil, err)
					equal(t, test.input, input)
				}

				srvJSON.WriteResponse(w, http.StatusOK, test.output)
			}))
			defer srv.Close()

			makeRequest(t, test.method, srv, test.input, test.output)
		})
	}
}
