package postgres

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/kcaashish/gotodo"
)

type UserStore struct {
	*sqlx.DB
}

func (s *UserStore) User(id uuid.UUID) (gotodo.User, error) {
	var u gotodo.User
	if err := s.Get(&u, `SELECT * FROM users WHERE id = $1`, id); err != nil {
		return gotodo.User{}, fmt.Errorf("Error getting User: %w", err)
	}
	return u, nil
}

func (s *UserStore) Users() ([]gotodo.User, error) {
	var uu []gotodo.User
	if err := s.Get(&uu, `SELECT * FROM users`); err != nil {
		return []gotodo.User{}, fmt.Errorf("Error getting Users: %w", err)
	}
	return uu, nil
}

func (s *UserStore) CreateUser(u *gotodo.User) error {
	if err := s.Get(u, `INSERT INTO users VALUES ($1, $2, $3, $4, $5, $6) RETURNING *`,
		u.ID, u.UserName, u.FirstName, u.LastName, u.Email, u.Password); err != nil {
		return fmt.Errorf("Error creating user: %w", err)
	}
	return nil
}

func (s *UserStore) UpdateUser(u *gotodo.User) error {
	if err := s.Get(u, `UPDATE users SET username = $2, first_name = $3, 
		last_name = $4, email = $5, password = $6 WHERE id = $1 RETURNING *`,
		u.ID, u.UserName, u.FirstName, u.LastName, u.Email, u.Password); err != nil {
		return fmt.Errorf("Error updating user: %w", err)
	}
	return nil
}

func (s *UserStore) DeleteUser(id uuid.UUID) error {
	if _, err := s.Exec(`DELETE FROM users WHERE id = $1`, id); err != nil {
		return fmt.Errorf("Error deleting user: %w", err)
	}
	return nil
}
