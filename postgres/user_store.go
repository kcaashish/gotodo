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
	if err := s.Select(&uu, `SELECT * FROM users`); err != nil {
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

func (s *UserStore) UpdateUser(id uuid.UUID, u *gotodo.User) error {
	if err := s.Get(u, `UPDATE users uu SET
		username = CASE WHEN $2 = '' THEN uu.username ELSE $2 END,
		first_name = CASE WHEN $3 = '' THEN uu.first_name ELSE $3 END,
		last_name = CASE WHEN $4 = '' THEN uu.last_name ELSE $4 END,
		email = CASE WHEN $5 = '' THEN uu.email ELSE $5 END,
		password = CASE WHEN $6 = '' THEN uu.password ELSE $6 END
		WHERE id = $1 RETURNING *`,
		id, u.UserName, u.FirstName, u.LastName, u.Email, u.Password); err != nil {
		return fmt.Errorf("Error updating user: %w", err)
	}
	return nil
}

func (s *UserStore) DeleteUser(id uuid.UUID) error {
	res, err := s.Exec(`DELETE FROM users WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("Error deleting user: %w", err)
	}

	if rows, _ := res.RowsAffected(); rows == 0 {
		return fmt.Errorf("Error deleting user: No such row")
	}
	return nil
}

func (s *UserStore) FindUser(email string) (gotodo.User, error) {
	var u gotodo.User
	er := s.Get(&u, `SELECT * FROM users WHERE email = $1`, email)
	if er != nil {
		return gotodo.User{}, fmt.Errorf("Invalid user! %w", er)
	}
	return u, nil
}
