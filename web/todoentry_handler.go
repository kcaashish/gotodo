package web

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/kcaashish/gotodo"
)

func (s *Server) getTodoEntry() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, _ := uuid.Parse(getField(r, 0))

		te, err := s.store.TodoEntry(id)
		// er if the given id is not in the database
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(te)
	}
}

func (s *Server) getTodoEntries() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tee, err := s.store.TodoEntriesByList()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tee)
	}
}

func (s *Server) createTodoEntry() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		te := &gotodo.TodoEntry{}
		if err := json.NewDecoder(r.Body).Decode(te); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if er := s.store.CreateTodoEntry(te); er != nil {
			http.Error(w, er.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(te)
	}
}

func (s *Server) updateTodoEntry() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, _ := uuid.Parse(getField(r, 0))
		te := &gotodo.TodoEntry{}
		if err := json.NewDecoder(r.Body).Decode(te); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if er := s.store.UpdateTodoEntry(id, te); er != nil {
			http.Error(w, er.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(te)
	}
}

func (s *Server) deleteTodoEntry() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, _ := uuid.Parse(getField(r, 0))

		if er := s.store.DeleteTodoEntry(id); er != nil {
			http.Error(w, er.Error(), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
