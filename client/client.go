package client

import (
	"bytes"
	"context"
	"io"
	"net/http"

	"github.com/gromey/proto-rest/coder"
	"github.com/gromey/proto-rest/logger"
)

type Client interface {
	coder.Coder
	Request(ctx context.Context, method, url string, body any, f func(*http.Request)) (*http.Response, error)
}

type protoClient struct {
	coder.Coder
	*http.Client
}

// New returns a new Client.
func New(coder coder.Coder, client *http.Client) Client {
	return &protoClient{Coder: coder, Client: client}
}

// Request sends an HTTP request based on the given method, URL, and optional body, and returns an HTTP response.
// To add additional data to the request, use the optional function f.
func (c protoClient) Request(ctx context.Context, method, url string, body any, f func(*http.Request)) (*http.Response, error) {
	var err error
	var reader io.Reader
	if body != nil {
		buf := new(bytes.Buffer)
		if err = c.Encode(buf, body); err != nil {
			return nil, err
		}
		reader = buf
	}

	var request *http.Request
	if request, err = http.NewRequestWithContext(ctx, method, url, reader); err != nil {
		return nil, err
	}

	if reader != nil {
		setContentType(request, c.ContentType())
	}

	if f != nil {
		f(request)
	}

	if logger.InLevel(logger.LevelTrace) && request != nil {
		logger.DumpHttpRequest(request, logger.Trace)
	}

	var response *http.Response
	if response, err = c.Do(request); err != nil {
		return nil, err
	}

	if logger.InLevel(logger.LevelTrace) && response != nil {
		logger.DumpHttpResponse(response, logger.Trace)
	}

	return response, nil
}

func setContentType(request *http.Request, contentType string) {
	if contentType != "" {
		request.Header.Set(coder.ContentType, contentType)
	}
}
