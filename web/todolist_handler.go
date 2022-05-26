package web

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/kcaashish/gotodo"
)

func (s *Server) getTodoList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := uuid.Parse(getField(r, 0))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

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

func (s *Server) createTodolist() http.HandlerFunc {
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

func (s *Server) updateTodolist() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		todolist := &gotodo.TodoList{}
		if err := json.NewDecoder(r.Body).Decode(todolist); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if er := s.store.UpdateTodoList(todolist); er != nil {
			http.Error(w, er.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(todolist)
	}
}

func (s *Server) deleteTodolist() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := uuid.Parse(getField(r, 0))

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if er := s.store.DeleteTodoList(id); er != nil {
			http.Error(w, er.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(id)
	}
}
