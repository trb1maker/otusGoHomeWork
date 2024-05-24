package internalhttp

import "net/http"

func (s *Server) ping(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("pong"))
}
