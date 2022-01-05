package server

import "net/http"

type Server struct {
	srv *http.Server
}

func New(addr string, h http.Handler) *Server {
	return &Server{
		srv: &http.Server{
			Addr:    addr,
			Handler: h,
		},
	}
}

func (s *Server) Start() error {
	return s.srv.ListenAndServe()
}
