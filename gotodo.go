package gotodo

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `db:"id"`
	UserName  string    `db:"username"`
	FirstName string    `db:"first_name"`
	LastName  string    `db:"last_name"`
	Email     string    `db:"email"`
	Password  string    `db:"password"`
}

type TodoList struct {
	ID          uuid.UUID `db:"id"`
	UserID      uuid.UUID `db:"user_id"`
	Title       string    `db:"title"`
	Description string    `db:"description"`
	CreatedDate time.Time `db:"created_date"`
	UpdatedDate time.Time `db:"updated_date"`
	DueDate     time.Time `db:"due_date"`
	Completed   bool      `db:"completed"`
}

type TodoEntry struct {
	ID          uuid.UUID `db:"id"`
	TodoListID  uuid.UUID `db:"todolist_id"`
	Content     string    `db:"content"`
	CreatedDate time.Time `db:"created_date"`
	UpdatedDate time.Time `db:"updated_date"`
	DueDate     time.Time `db:"due_date"`
	Completed   bool      `db:"completed"`
}

type TodoListStore interface {
	TodoList(id uuid.UUID) (TodoList, error)
	TodoLists() ([]TodoList, error)
	CreateTodoList(t *TodoList) error
	UpdateTodoList(t *TodoList) error
	DeleteTodoList(id uuid.UUID) error
}

type TodoEntryStore interface {
	TodoEntry(id uuid.UUID) (TodoEntry, error)
	TodoEntriesByList(todolist_id uuid.UUID) ([]TodoEntry, error)
	CreateTodoEntry(t *TodoEntry) error
	UpdateTodoEntry(t *TodoEntry) error
	DeleteTodoEntry(id uuid.UUID) error
}
