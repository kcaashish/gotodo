package web

import "github.com/kcaashish/gotodo"

func NewServer(store gotodo.Store) *Server {
	s := &Server{
		store: store,
	}
	return s
}

type Server struct {
	store gotodo.Store
}
