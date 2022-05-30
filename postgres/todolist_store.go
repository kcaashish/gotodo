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

func (s *TodoListStore) UpdateTodoList(id uuid.UUID, t *gotodo.TodoList) error {
	if err := s.Get(t, `UPDATE todo_list tl SET 
		user_id = CASE WHEN NULLIF($2, '00000000-0000-0000-0000-000000000000'::UUID) IS NULL THEN tl.user_id ELSE $2::UUID END, 
		title = CASE WHEN $3 = '' THEN tl.title ELSE $3 END, 
		description = CASE WHEN $4 = '' THEN tl.description ELSE $4 END, 
		created_date = CASE WHEN NULLIF($5, '0001-01-01T00:00:00Z'::TIMESTAMP) IS NULL THEN tl.created_date ELSE $5::TIMESTAMP END, 
		updated_date = CASE WHEN NULLIF($6,'0001-01-01T00:00:00Z'::TIMESTAMP) IS NULL THEN tl.updated_date ELSE $6::TIMESTAMP END, 
		due_date = CASE WHEN NULLIF($7,'0001-01-01T00:00:00Z'::TIMESTAMP) IS NULL THEN tl.due_date ELSE $7::TIMESTAMP END, 
		completed = CASE WHEN NULLIF($8,'') IS NULL THEN tl.completed ELSE $8::BOOLEAN END 
		WHERE id = $1 RETURNING *`,
		id, t.UserID, t.Title, t.Description, t.CreatedDate,
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
