package gotodo

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `db:"id" json:"id"`
	UserName  string    `db:"username" json:"username"`
	FirstName string    `db:"first_name" json:"first_name"`
	LastName  string    `db:"last_name" json:"last_name"`
	Email     string    `db:"email" json:"email"`
	Password  string    `db:"password" json:"password"`
}

type TodoList struct {
	ID          uuid.UUID `db:"id" json:"id"`
	UserID      uuid.UUID `db:"user_id" json:"user_id"`
	Title       string    `db:"title" json:"title"`
	Description string    `db:"description" json:"description"`
	CreatedDate time.Time `db:"created_date" json:"creaded_date"`
	UpdatedDate time.Time `db:"updated_date" json:"updated_date"`
	DueDate     time.Time `db:"due_date" json:"due_date"`
	Completed   bool      `db:"completed" json:"completed"`
}

type TodoEntry struct {
	ID          uuid.UUID `db:"id" json:"id"`
	TodoListID  uuid.UUID `db:"todolist_id" json:"todolist_id"`
	Content     string    `db:"content" json:"content"`
	CreatedDate time.Time `db:"created_date" json:"created_date"`
	UpdatedDate time.Time `db:"updated_date" json:"updated_date"`
	DueDate     time.Time `db:"due_date" json:"due_date"`
	Completed   bool      `db:"completed" json:"completed"`
}

type UserStore interface {
	User(id uuid.UUID) (User, error)
	Users() ([]User, error)
	CreateUser(u *User) error
	UpdateUser(u *User) error
	DeleteUser(id uuid.UUID) error
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
	TodoEntriesByList() ([]TodoEntry, error)
	CreateTodoEntry(t *TodoEntry) error
	UpdateTodoEntry(t *TodoEntry) error
	DeleteTodoEntry(id uuid.UUID) error
}

type Store interface {
	UserStore
	TodoListStore
	TodoEntryStore
}
