package server

import "net/http"

type Server struct {
	srv *http.Server
}

func New(s *http.Server) *Server {
	return &Server{
		srv: s,
	}
}

func (s *Server) Start() error {
	return s.srv.ListenAndServe()
}
