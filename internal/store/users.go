package store

import (
	"context"
	"database/sql"
)

type User struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"-"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type UserStore struct {
	db *sql.DB
}

func (s *UserStore) Create(ctx context.Context, user *User) error {
	query := `
	INSERT INTO users (username, email, password)
	VALUES (?, ?, ?) RETURNING id, created_at, updated_at;`

	err := s.db.QueryRowContext(
		ctx,
		query,
		user.Username,
		user.Email,
		user.Password,
	).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func (s *UserStore) GetByIDOrEmail(ctx context.Context, id int64, email string) (*User, error) {
	query := `
	SELECT id, username, email, created_at, updated_at
	FROM users
	WHERE id = ? OR email = ?;`
	// Initialize an empty User struct to hold the result
	user := &User{}

	// Execute the query and scan the result into the user struct
	err := s.db.QueryRowContext(ctx, query, id, email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserStore) Update(ctx context.Context, user *User) error {
	query := `
	UPDATE users
	SET username = ?, email = ?, password = ?, updated_at = CURRENT_TIMESTAMP
	WHERE id = ?;`

	_, err := s.db.ExecContext(
		ctx,
		query,
		user.Username,
		user.Email,
		user.Password,
		user.ID,
	)

	if err != nil {
		return err
	}

	return nil
}
