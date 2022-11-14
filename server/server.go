package server

import (
	"net/http"

	"github.com/gromey/proto-rest/coder"
	"github.com/gromey/proto-rest/logger"
)

type Server interface {
	coder.Coder
	WriteResponse(w http.ResponseWriter, statusCode int, v any)
}

type protoServer struct {
	coder.Coder
}

// New returns a new Server.
func New(coder coder.Coder) Server {
	return &protoServer{Coder: coder}
}

// WriteResponse encodes the value pointed to by v and writes it and statusCode to the stream.
func (s *protoServer) WriteResponse(w http.ResponseWriter, statusCode int, v any) {
	if v != nil {
		if w.Header().Get(coder.ContentType) == "" {
			if t := s.ContentType(); t != "" {
				w.Header().Set(coder.ContentType, t)
			}
		}
		w.WriteHeader(statusCode)
		if err := s.Encode(w, v); err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			if logger.InLevel(logger.LevelError) {
				logger.Error("Can't encode response. Error: ", err)
			}
		}
		return
	}
	w.WriteHeader(statusCode)
}
