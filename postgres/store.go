package postgres

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func NewStore(dataSourceName string) (*Store, error) {
	db, err := sqlx.Open("postgres", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("Error opening database: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("Error connecting to the database: %w", err)
	}

	return &Store{
		UserStore:      &UserStore{DB: db},
		TodoListStore:  &TodoListStore{DB: db},
		TodoEntryStore: &TodoEntryStore{DB: db},
	}, nil
}

type Store struct {
	*UserStore
	*TodoListStore
	*TodoEntryStore
}
