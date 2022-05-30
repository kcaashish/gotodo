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

func (s *TodoEntryStore) TodoEntriesByList() ([]gotodo.TodoEntry, error) {
	var tee []gotodo.TodoEntry
	if err := s.Select(&tee, `SELECT * FROM todo_entry`); err != nil {
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

func (s *TodoEntryStore) UpdateTodoEntry(id uuid.UUID, te *gotodo.TodoEntry) error {
	if err := s.Get(te, `UPDATE todo_entry t SET 
		todolist_id = CASE WHEN NULLIF($2, '00000000-0000-0000-0000-000000000000'::UUID) IS NULL THEN t.todolist_id ELSE $2::uuid END, 
		content = CASE WHEN $3 = '' THEN t.content ELSE $3 END,
		created_date = CASE WHEN NULLIF($4, '0001-01-01T00:00:00Z'::TIMESTAMP) IS NULL THEN t.created_date ELSE $4::timestamp END, 
		updated_date = CASE WHEN NULLIF($5, '0001-01-01T00:00:00Z'::TIMESTAMP) IS NULL THEN t.updated_date ELSE $5::timestamp END, 
		due_date = CASE WHEN NULLIF($6, '0001-01-01T00:00:00Z'::TIMESTAMP) IS NULL THEN t.due_date ELSE $6::timestamp END, 
		completed = CASE WHEN NULLIF($7, '') IS NULL THEN t.completed ELSE $7::boolean END 
		WHERE id = $1 RETURNING *`,
		id, te.TodoListID, te.Content, te.CreatedDate,
		te.UpdatedDate, te.DueDate, te.Completed); err != nil {
		return fmt.Errorf("Error updating TodoEntry: %w", err)
	}
	return nil
}

func (s *TodoEntryStore) DeleteTodoEntry(id uuid.UUID) error {
	res, err := s.Exec(`DELETE FROM todo_entry WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("Error in deleting TodoEntry: %w", err)
	}

	if rows, _ := res.RowsAffected(); rows == 0 {
		return fmt.Errorf("Error deleting user: No such row")
	}
	return nil
}
