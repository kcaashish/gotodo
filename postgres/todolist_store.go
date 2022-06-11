package postgres

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/kcaashish/gotodo"
)

type TodoListStore struct {
	*sqlx.DB
}

func (s *TodoListStore) TodoList(id uuid.UUID) (gotodo.TodoList, error) {
	var t gotodo.TodoList
	if err := s.Get(&t, `SELECT * FROM todo_list WHERE id = $1`, id); err != nil {
		return gotodo.TodoList{}, fmt.Errorf("Error getting TodoList: %w", err)
	}
	return t, nil
}

func (s *TodoListStore) TodoLists() ([]gotodo.TodoList, error) {
	var tl []gotodo.TodoList
	if err := s.Select(&tl, `SELECT * FROM todo_list`); err != nil {
		return []gotodo.TodoList{}, fmt.Errorf("Error getting TodoList: %w", err)
	}
	return tl, nil
}

func (s *TodoListStore) CreateTodoList(t *gotodo.TodoList) error {
	if err := s.Get(t, `INSERT INTO todo_list(id, user_id, title, description, due_at, completed) VALUES ($1, $2, $3, $4, $5, $6) RETURNING *`,
		t.ID, t.UserID, t.Title, t.Description, t.DueAt, t.Completed); err != nil {
		return fmt.Errorf("Error creating TodoList: %w", err)
	}
	return nil
}

func (s *TodoListStore) UpdateTodoList(id uuid.UUID, t *gotodo.TodoList) error {
	if err := s.Get(t, `UPDATE todo_list tl SET 
		title = CASE WHEN $2 = '' THEN tl.title ELSE $2 END, 
		description = CASE WHEN $3 = '' THEN tl.description ELSE $3 END, 
		updated_at = $4::TIMESTAMP, 
		due_at = CASE WHEN NULLIF($5,'0001-01-01T00:00:00Z'::TIMESTAMP) IS NULL THEN tl.due_at ELSE $5::TIMESTAMP END, 
		completed = CASE WHEN NULLIF($6,'') IS NULL THEN tl.completed ELSE $6::BOOLEAN END 
		WHERE id = $1 RETURNING *`,
		id, t.Title, t.Description, t.UpdatedAt.Time, t.DueAt, t.Completed); err != nil {
		return fmt.Errorf("Error updating TodoList: %w", err)
	}
	return nil
}

func (s *TodoListStore) DeleteTodoList(id uuid.UUID) error {
	res, err := s.Exec(`DELETE FROM todo_list WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("Error in deleting TodoList: %w", err)
	}

	if rows, _ := res.RowsAffected(); rows == 0 {
		return fmt.Errorf("Error deleting user: No such row")
	}
	return nil
}
