package web

import (
	"fmt"
	"net/http"
)

func (s *Server) getTodoLists() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tl, err := s.store.TodoLists()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Fprint(w, tl)
	}
}
