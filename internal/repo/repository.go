// Package repo provides the repository layer abstraction for data access operations.
// It implements the repository pattern with interface-based design for dependency injection
// and testability across different storage backends.
package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// Common repository errors
var (
	// ErrRepositoryNotInitialized is returned when repository is not properly initialized
	ErrRepositoryNotInitialized = errors.New("repository not initialized")

	// ErrInvalidInput is returned when invalid input is provided to repository methods
	ErrInvalidInput = errors.New("invalid input provided")
)

// PostsRepository defines the contract for post-related database operations.
// This interface allows for easy testing and potential future implementations
// with different storage backends (Redis, MongoDB, etc.).
type PostsRepository interface {
	Create(ctx context.Context, post *Post) error
	GetByID(ctx context.Context, id int64) (*Post, error)
}

type UsersRepository interface {
	Create(ctx context.Context, user *User) error
}

type CommentsRepository interface {
	GetByPostID(ctx context.Context, postID int64) ([]Comment, error)
}

// Repository aggregates all repository interfaces into a single structure.
// This provides a unified access point for all data operations and simplifies
// dependency injection in the service layer.
//
// Usage:
//
//	repo := NewRepository(db)
//	post, err := repo.Posts.GetByID(ctx, 123)
//	user, err := repo.Users.GetByEmail(ctx, "user@example.com")
type Repository struct {
	Posts    PostsRepository
	Users    UsersRepository
	Comments CommentsRepository
}

// NewRepository creates a new Repository instance with PostgreSQL implementations.
// It initializes all repository interfaces with concrete PostgreSQL store implementations
// and validates the database connection.
//
// The function will:
//   - Validate the database connection is not nil
//   - Initialize PostStore and UserStore with the database connection
//   - Return a fully configured Repository ready for use
//
// Example:
//
//	db, err := sql.Open("postgres", dsn)
//	if err != nil {
//	  log.Fatal(err)
//	}
//
//	repo := repo.NewRepository(db)
//	post, err := repo.Posts.GetByID(ctx, postID)
func NewPostgresRepo(db *sql.DB) (*Repository, error) {
	if db == nil {
		return nil, errors.New("database connection cannot be nil")
	}

	var store UsersRepository
	store = &UserStore{db}

	fmt.Println(store)

	return &Repository{
		Posts:    NewPostStore(db),
		Users:    &UserStore{db},
		Comments: &CommentRepo{db},
	}, nil

}

// Health performs a health check on all repository components.
// It verifies database connectivity and returns any errors encountered.
func (r *Repository) Health(ctx context.Context) error {
	if r.Posts == nil || r.Users == nil {
		return ErrRepositoryNotInitialized
	}

	// If stores implement health check interface, call it
	if healthChecker, ok := r.Posts.(interface{ Health(context.Context) error }); ok {
		if err := healthChecker.Health(ctx); err != nil {
			return err
		}
	}

	if healthChecker, ok := r.Users.(interface{ Health(context.Context) error }); ok {
		if err := healthChecker.Health(ctx); err != nil {
			return err
		}
	}

	return nil
}
