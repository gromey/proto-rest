package middleware

import (
	"net/http"
	"strconv"
	"strings"
)

// CORSOptions represents a functional option for configuring the CORS middleware.
type CORSOptions struct {
	AllowedOrigins   []string // The origins that the server allows.
	AllowMethods     []string // List of methods that the server allows.
	AllowHeaders     []string // List of headers that the server allows.
	MaxAge           int      // Tells the browser how long (in seconds) to cache the response to the preflight request.
	AllowCredentials bool     // Allow browsers to expose the response to the external JavaScript code.
}

// AllowCORS sets headers for CORS mechanism supports secure.
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
