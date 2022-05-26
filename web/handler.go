package web

import (
	"encoding/json"
	"net/http"
)

func (s *Server) getTodoLists() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tl, err := s.store.TodoLists()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tl)
	}
}
