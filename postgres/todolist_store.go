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
	if err := s.Get(t, `INSERT INTO todo_list VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING *`,
		t.ID, t.UserID, t.Title, t.Description, t.CreatedDate, t.UpdatedDate, t.DueDate, t.Completed); err != nil {
		return fmt.Errorf("Error creating TodoList: %w", err)
	}
	return nil
}

func (s *TodoListStore) UpdateTodoList(t *gotodo.TodoList) error {
	if err := s.Get(t,
		`UPDATE todo_list SET user_id = $2, title = $3, description = $4, 
		created_date = $5, updated_date = $6, due_date = $7, 
		completed = $8 WHERE id = $1 RETURNING *`,
		t.ID, t.UserID, t.Title, t.Description, t.CreatedDate,
		t.UpdatedDate, t.DueDate, t.Completed); err != nil {
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
