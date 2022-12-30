package roundtripper

import (
	"net/http"
	"runtime/debug"
	"time"

	"github.com/gromey/proto-rest/logger"
)

type function func(*http.Request) (*http.Response, error)

func (f function) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

// Sequencer chains http.RoundTrippers in a chain.
func Sequencer(rts ...func(http.RoundTripper) http.RoundTripper) http.RoundTripper {
	rt := http.DefaultTransport
	for _, f := range rts {
		rt = f(rt)
	}
	return rt
}

// Timer measures the time taken by http.RoundTripper.
func Timer(next http.RoundTripper) http.RoundTripper {
	return function(func(r *http.Request) (*http.Response, error) {
		defer func(start time.Time) {
			if logger.InLevel(logger.LevelDebug) {
				logger.Debugf("%s %s %s", r.Method, r.RequestURI, time.Since(start))
			}
		}(time.Now())
		return next.RoundTrip(r)
	})
}

// PanicCatcher handles panics in http.RoundTripper.
func PanicCatcher(next http.RoundTripper) http.RoundTripper {
	return function(func(r *http.Request) (*http.Response, error) {
		defer func() {
			if rec := recover(); rec != nil {
				if logger.InLevel(logger.LevelError) {
					logger.Error(string(debug.Stack()))
				}
			}
		}()
		return next.RoundTrip(r)
	})
}

// DumpHttp dumps the HTTP request and response, and prints out with logFunc.
func DumpHttp(logLevel logger.Level, logFunc func(v ...any)) func(next http.RoundTripper) http.RoundTripper {
	return func(next http.RoundTripper) http.RoundTripper {
		return function(func(r *http.Request) (*http.Response, error) {
			if logger.InLevel(logLevel) {
				logger.DumpHttpRequest(r, logFunc)

				resp, err := next.RoundTrip(r)
				if err != nil {
					return nil, err
				}

				logger.DumpHttpResponse(resp, logFunc)

				return resp, nil
			}

			return next.RoundTrip(r)
		})
	}
}
