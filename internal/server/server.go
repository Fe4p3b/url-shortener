package server

import "net/http"

type Server struct {
	srv *http.Server
}

func NewServer(s *http.Server) *Server {
	return &Server{
		srv: s,
	}
}

func (s *Server) Start(addr string) error {
	s.srv.Addr = addr
	return s.srv.ListenAndServe()
}
