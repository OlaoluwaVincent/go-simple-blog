package store

import (
	"context"
	"database/sql"
)

type UserInterface interface {
	Create(context.Context, *User) (*User, error)
	GetByIDUsernameOrEmail(ctx context.Context, id *int64, username *string, email *string) (*User, error)
	Update(context.Context, *User) (*User, error)
}

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

func (s *UserStore) Create(ctx context.Context, user *User) (*User, error) {
	query := `
	INSERT INTO users (username, email, password)
	VALUES ($1, $2, $3)
	RETURNING id, username, email, created_at, updated_at;
	`

	err := s.db.QueryRowContext(
		ctx,
		query,
		user.Username,
		user.Email,
		user.Password,
	).Scan(
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

func (s *UserStore) GetByIDUsernameOrEmail(
	ctx context.Context,
	id *int64,
	username *string,
	email *string,
) (*User, error) {

	query := `
	SELECT id, username, email, password, created_at, updated_at
	FROM users
	WHERE
		($1::bigint IS NOT NULL AND id = $1)
		OR ($2::text IS NOT NULL AND username = $2)
		OR ($3::text IS NOT NULL AND email = $3)
	LIMIT 1;
	`

	user := &User{}

	err := s.db.QueryRowContext(ctx, query, id, username, email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserStore) Update(ctx context.Context, user *User) (*User, error) {
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
		return nil, err
	}

	return user, nil
}
