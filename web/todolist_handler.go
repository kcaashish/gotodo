package web

import (
	"encoding/json"
	"net/http"
	"time"

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

		// fetch from db first then check for UserID
		userid := r.Context().Value("user").(uuid.UUID)
		if userid != t.UserID {
			http.Error(w, "Invalid request", http.StatusForbidden)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(t)
	}
}

func (s *Server) getTodoLists() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		todolists, err := s.store.TodoLists()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// return back the tls of particular user only
		userid := r.Context().Value("user").(uuid.UUID)
		tlToReturn := make([]gotodo.TodoList, 0)
		for _, tl := range todolists {
			if userid == tl.UserID {
				tlToReturn = append(tlToReturn, tl)
			}
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tlToReturn)
	}
}

func (s *Server) createTodoList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		todolist := &gotodo.TodoList{}
		todolist.ID = uuid.New()

		userid := r.Context().Value("user").(uuid.UUID)
		todolist.UserID = userid

		todolist.CreatedDate = time.Now().Local()
		todolist.UpdatedDate = time.Now().Local()
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
		todolist.UpdatedDate = time.Now().Local()
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

		// first fetch the tl based on id
		t, er := s.store.TodoList(id)
		if er != nil {
			http.Error(w, er.Error(), http.StatusInternalServerError)
		}

		// fetch from db first then check for UserID
		userid := r.Context().Value("user").(uuid.UUID)
		if userid != t.UserID {
			http.Error(w, "Invalid request", http.StatusForbidden)
			return
		}

		if er := s.store.DeleteTodoList(id); er != nil {
			http.Error(w, er.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(t)
	}
}
