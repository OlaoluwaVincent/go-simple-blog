package store

import (
	"database/sql"
)

type Storage struct {
	Posts PostInterface
	Users UserInterface
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Posts: &PostStore{db: db},
		Users: &UserStore{db: db},
	}
}
