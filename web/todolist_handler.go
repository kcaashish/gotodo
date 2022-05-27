package web

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/kcaashish/gotodo"
)

func (s *Server) getTodoList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, _ := uuid.Parse(getField(r, 0))

		t, er := s.store.TodoList(id)
		if er != nil {
			http.Error(w, er.Error(), http.StatusInternalServerError)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(t)
	}
}

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

func (s *Server) createTodoList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		todolist := &gotodo.TodoList{}
		if err := json.NewDecoder(r.Body).Decode(todolist); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if er := s.store.CreateTodoList(todolist); er != nil {
			http.Error(w, er.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(todolist)
	}
}

func (s *Server) updateTodoList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, _ := uuid.Parse(getField(r, 0))
		todolist := &gotodo.TodoList{}
		if err := json.NewDecoder(r.Body).Decode(todolist); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if er := s.store.UpdateTodoList(id, todolist); er != nil {
			http.Error(w, er.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(todolist)
	}
}

func (s *Server) deleteTodoList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, _ := uuid.Parse(getField(r, 0))

		if er := s.store.DeleteTodoList(id); er != nil {
			http.Error(w, er.Error(), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
