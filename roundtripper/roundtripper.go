package roundtripper

import (
	"net/http"
	"runtime/debug"
	"time"

	"github.com/gromey/proto-rest/logger"
)

// The Func type is an adapter to allow the use of ordinary functions as HTTP round trippers.
// If f is a function with the appropriate signature, Func(f) is a RoundTripper that calls f.
type Func func(*http.Request) (*http.Response, error)

// RoundTrip calls f(r).
func (f Func) RoundTrip(r *http.Request) (*http.Response, error) {
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
func Timer(logLevel logger.Level) func(http.RoundTripper) http.RoundTripper {
	return func(next http.RoundTripper) http.RoundTripper {
		return Func(func(r *http.Request) (*http.Response, error) {
			defer func(start time.Time) {
				if logger.InLevel(logLevel) {
					logLevel.Printf()("%s %s %s", r.Method, r.RequestURI, time.Since(start))
				}
			}(time.Now())
			return next.RoundTrip(r)
		})
	}
}

// PanicCatcher handles panics in http.RoundTripper.
func PanicCatcher(next http.RoundTripper) http.RoundTripper {
	return Func(func(r *http.Request) (*http.Response, error) {
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
func DumpHttp(logLevel logger.Level) func(http.RoundTripper) http.RoundTripper {
	return func(next http.RoundTripper) http.RoundTripper {
		return Func(func(r *http.Request) (*http.Response, error) {
			if logger.InLevel(logLevel) {
				logger.DumpHttpRequest(r, logLevel.Print())

				resp, err := next.RoundTrip(r)
				if err != nil {
					return nil, err
				}

				logger.DumpHttpResponse(resp, logLevel.Print())

				return resp, nil
			}

			return next.RoundTrip(r)
		})
	}
}
