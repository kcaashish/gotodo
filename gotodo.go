package gotodo

import (
	"database/sql"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type User struct {
	ID                uuid.UUID `db:"id" json:"id"`
	UserName          string    `db:"username" json:"username"`
	FirstName         string    `db:"first_name" json:"first_name"`
	LastName          string    `db:"last_name" json:"last_name"`
	Email             string    `db:"email" json:"email"`
	Password          string    `db:"password" json:"password"`
	CreatedAt         time.Time `db:"created_at" json:"created_at"`
	PasswordChangedAt time.Time `db:"password_changed_at" json:"password_changed_at"`
}

type UserResponse struct {
	ID                uuid.UUID `json:"id"`
	UserName          string    `json:"username"`
	FirstName         string    `json:"first_name"`
	LastName          string    `json:"last_name"`
	Email             string    `json:"email"`
	CreatedAt         time.Time `json:"created_at"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
}

type LoginUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginUserResponse struct {
	Token Token        `json:"token"`
	User  UserResponse `json:"user"`
}

type TodoList struct {
	ID          uuid.UUID    `db:"id" json:"id"`
	UserID      uuid.UUID    `db:"user_id" json:"user_id"`
	Title       string       `db:"title" json:"title"`
	Description string       `db:"description" json:"description"`
	CreatedAt   time.Time    `db:"created_at" json:"creaded_at"`
	UpdatedAt   sql.NullTime `db:"updated_at" json:"updated_at"`
	DueAt       time.Time    `db:"due_at" json:"due_at"`
	Completed   bool         `db:"completed" json:"completed"`
}

type TodoEntry struct {
	ID         uuid.UUID    `db:"id" json:"id"`
	TodoListID uuid.UUID    `db:"todolist_id" json:"todolist_id"`
	Content    string       `db:"content" json:"content"`
	CreatedAt  time.Time    `db:"created_at" json:"created_at"`
	UpdatedAt  sql.NullTime `db:"updated_at" json:"updated_at"`
	DueAt      time.Time    `db:"due_at" json:"due_at"`
	Completed  bool         `db:"completed" json:"completed"`
}

type Claims struct {
	UserID   uuid.UUID `json:"user_id"`
	UserName string    `json:"user_name"`
	Email    string    `json:"email"`
	*jwt.StandardClaims
}

type Token struct {
	AccessToken           string    `json:"access_token"`
	AccessTokenExpiresAt  time.Time `json:"access_token_expires_at"`
	RefreshToken          string    `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time `json:"refresh_token_expires_at"`
}

type UserStore interface {
	User(id uuid.UUID) (User, error)
	Users() ([]User, error)
	CreateUser(u *User) error
	UpdateUser(id uuid.UUID, u *User) error
	DeleteUser(id uuid.UUID) error
	FindUser(email string) (User, error)
}

type TodoListStore interface {
	TodoList(id uuid.UUID) (TodoList, error)
	TodoLists() ([]TodoList, error)
	CreateTodoList(t *TodoList) error
	UpdateTodoList(id uuid.UUID, t *TodoList) error
	DeleteTodoList(id uuid.UUID) error
}

type TodoEntryStore interface {
	TodoEntry(id uuid.UUID) (TodoEntry, error)
	TodoEntriesByList() ([]TodoEntry, error)
	CreateTodoEntry(t *TodoEntry) error
	UpdateTodoEntry(id uuid.UUID, t *TodoEntry) error
	DeleteTodoEntry(id uuid.UUID) error
}

type Store interface {
	UserStore
	TodoListStore
	TodoEntryStore
}
