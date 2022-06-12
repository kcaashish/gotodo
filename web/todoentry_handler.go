package web

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/kcaashish/gotodo"
)

func (s *Server) getTodoEntry() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tl_id, _ := uuid.Parse(getField(r, 0))
		id, _ := uuid.Parse(getField(r, 1))

		te, err := s.store.TodoEntry(id)
		// er if the given id is not in the database
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if te.TodoListID != tl_id {
			http.Error(w, "Invalid request", http.StatusForbidden)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(te)
	}
}

func (s *Server) getTodoEntries() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tl_id, _ := uuid.Parse(getField(r, 0))
		tee, err := s.store.TodoEntriesByList()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		teToReturn := make([]gotodo.TodoEntry, 0)
		for _, te := range tee {
			if te.TodoListID == tl_id {
				teToReturn = append(teToReturn, te)
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(teToReturn)
	}
}

func (s *Server) createTodoEntry() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tl_id, _ := uuid.Parse(getField(r, 0))

		te := &gotodo.TodoEntry{}
		te.ID = uuid.New()
		te.TodoListID = tl_id

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
		tl_id, _ := uuid.Parse(getField(r, 0))
		id, _ := uuid.Parse(getField(r, 1))

		// first fetch te based on id
		te_check, err := s.store.TodoEntry(id)
		// er if the given id is not in the database
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if te_check.TodoListID != tl_id {
			http.Error(w, "Invalid request", http.StatusForbidden)
			return
		}

		// if valid, proceed with update
		te := &gotodo.TodoEntry{}

		te.UpdatedAt = sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		}

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
		tl_id, _ := uuid.Parse(getField(r, 0))
		id, _ := uuid.Parse(getField(r, 1))

		// first fetch te based on id
		te_check, err := s.store.TodoEntry(id)
		// er if the given id is not in the database
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if te_check.TodoListID != tl_id {
			http.Error(w, "Invalid request", http.StatusForbidden)
			return
		}

		if er := s.store.DeleteTodoEntry(id); er != nil {
			http.Error(w, er.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(te_check)
	}
}
