package middleware

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	"github.com/gromey/proto-rest/logger"
)

// Timer middleware measures the time taken by http.HandlerFunc.
func Timer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		if logger.InLevel(logger.LevelDebug) {
			logger.Debugf("%s %s %s", r.Method, r.RequestURI, time.Since(start))
		}
	})
}

// PanicCatcher middleware handles panics in http.HandlerFunc.
func PanicCatcher(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				if logger.InLevel(logger.LevelError) {
					logger.Error(string(debug.Stack()))
				}
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// CORSOptions represents a functional option for configuring the CORS middleware.
type CORSOptions struct {
	AllowedOrigins   []string // The origins that the server allows.
	AllowMethods     []string // List of methods that the server allows.
	AllowHeaders     []string // List of headers that the server allows.
	MaxAge           int      // Tells the browser how long (in seconds) to cache the response to the preflight request.
	AllowCredentials bool     // Allow browsers to expose the response to the external JavaScript code.
}

// AllowCORS middleware sets headers for CORS mechanism supports secure.
func AllowCORS(next http.Handler, opts *CORSOptions) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if origin := r.Header.Get("Origin"); checkOrigin(origin, opts.AllowedOrigins) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}
		if r.Method == http.MethodOptions {
			w.Header().Set("Access-Control-Allow-Methods", strings.Join(opts.AllowMethods, ","))
			w.Header().Set("Access-Control-Allow-Headers", strings.Join(opts.AllowHeaders, ","))
			w.Header().Set("Access-Control-Max-Age", strconv.Itoa(opts.MaxAge))
			if opts.AllowCredentials {
				w.Header().Set("Access-Control-Allow-Credentials", "true")
			}
			return
		}
		next.ServeHTTP(w, r)
	})
}

func checkOrigin(origin string, allowedOrigins []string) bool {
	if origin == "" {
		return false
	}
	for _, v := range allowedOrigins {
		if origin == v || v == "*" {
			return true
		}
	}
	return false
}

// DumpHttp dumps the HTTP request and response, and prints out with logFunc.
func DumpHttp(logLevel logger.Level, logFunc func(v ...any)) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if logger.InLevel(logLevel) {
				logger.DumpHttpRequest(r, logFunc)

				buf := new(bytes.Buffer)
				recorder := httptest.NewRecorder()

				next.ServeHTTP(recorder, r)

				for key, values := range recorder.Header() {
					w.Header().Del(key)
					for _, value := range values {
						w.Header().Set(key, value)
					}
				}

				_, _ = recorder.Body.WriteTo(io.MultiWriter(w, buf))
				recorder.Body = buf

				logger.DumpHttpResponse(recorder.Result(), logFunc)

				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
