package repo

import (
	"context"
	"database/sql"
)

type UserStore struct {
	db *sql.DB
}

// User represents a user in the system
//
//	@Description	User account information
type User struct {
	// User's unique identifier
	//	@example	1
	ID int64 `json:"id" example:"1"`

	// User's unique username
	//	@example	john_doe
	Username string `json:"username" example:"john_doe"`

	// User's email address
	//	@example	john.doe@example.com
	Email string `json:"email" example:"john.doe@example.com"`

	// User's password (never returned in responses)
	Password string `json:"-"`

	// Timestamp when the user account was created
	//	@example	2024-01-15T10:30:00Z
	CreatedAt string `json:"created_at" example:"2024-01-15T10:30:00Z"`
}

func (s *UserStore) Create(ctx context.Context, user *User) error {
	query := `
		INSERT INTO users (username, password, email) VALUES ($1, $2, $3)
    	RETURNING id, created_at
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()

	//role := user.Role.Name
	//if role == "" {
	//	role = "user"
	//}

	err := s.db.QueryRowContext(
		ctx,
		query,
		user.Username,
		user.Password,
		user.Email,
	).Scan(
		&user.ID,
		&user.CreatedAt,
	)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserStore) GetByID(ctx context.Context, userID int64) (*User, error) {
	query := `
		SELECT users.id, username, email, password, created_at
		FROM users
		WHERE users.id = $1
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()

	user := &User{}
	err := s.db.QueryRowContext(
		ctx,
		query,
		userID,
	).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
	)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrPostNotFound
		default:
			return nil, err
		}
	}
	return user, nil
}
