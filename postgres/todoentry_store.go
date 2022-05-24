package postgres

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/kcaashish/gotodo"
)

type TodoEntryStore struct {
	*sqlx.DB
}

func (s *TodoEntryStore) TodoEntry(id uuid.UUID) (gotodo.TodoEntry, error) {
	var te gotodo.TodoEntry
	if err := s.Get(&te, `SELECT * FROM todo_entry WHERE id = $1`, id); err != nil {
		return gotodo.TodoEntry{}, fmt.Errorf("Error getting TodoEntry: %w", err)
	}
	return te, nil
}

func (s *TodoEntryStore) TodoEntriesByList(todolist_id uuid.UUID) ([]gotodo.TodoEntry, error) {
	var tee []gotodo.TodoEntry
	if err := s.Get(&tee, `SELECT * FROM todo_entry`); err != nil {
		return []gotodo.TodoEntry{}, fmt.Errorf("Error getting TodoEntry: %w", err)
	}
	return tee, nil
}

func (s *TodoEntryStore) CreateTodoEntry(te *gotodo.TodoEntry) error {
	if err := s.Get(te, `INSERT INTO todo_entry VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING *`,
		te.ID, te.TodoListID, te.Content, te.CreatedDate, te.UpdatedDate, te.DueDate, te.Completed); err != nil {
		return fmt.Errorf("Error creating TodoEntry: %w", err)
	}
	return nil
}

func (s *TodoEntryStore) UpdateTodoEntry(te *gotodo.TodoEntry) error {
	if err := s.Get(te,
		`UPDATE todo_entry SET todolist_id = $2, content = $3,
		created_date = $4, updated_date = $5, due_date = $6, 
		completed = $7 WHERE id = $1 RETURNING *`,
		te.ID, te.TodoListID, te.Content, te.CreatedDate,
		te.UpdatedDate, te.DueDate, te.Completed); err != nil {
		return fmt.Errorf("Error updating TodoEntry: %w", err)
	}
	return nil
}

func (s *TodoEntryStore) DeleteTodoEntry(id uuid.UUID) error {
	if _, err := s.Exec(`DELETE FROM todo_entry WHERE id = $1`, id); err != nil {
		return fmt.Errorf("Error in deleting TodoEntry: %w", err)
	}
	return nil
}
