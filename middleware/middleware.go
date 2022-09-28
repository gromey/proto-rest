package middleware

import (
	"net/http"
	"runtime/debug"
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
