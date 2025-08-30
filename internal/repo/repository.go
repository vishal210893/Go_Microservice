// Package repo provides the repository layer abstraction for data access operations.
// It implements the repository pattern with interface-based design for dependency injection
// and testability across different storage backends.
package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

// Common repository errors
var (
	ErrRepositoryNotInitialized = errors.New("repository not initialized")
	ErrInvalidInput             = errors.New("invalid input provided")
	QueryTimeout                = 15 * time.Second
	ErrConflict                 = errors.New("conflict error")
)

// PostsRepository defines the contract for post-related database operations.
// This interface allows for easy testing and potential future implementations
// with different storage backends (Redis, MongoDB, etc.).
type PostsRepository interface {
	Create(ctx context.Context, post *Post) error
	GetByID(ctx context.Context, id int64) (*Post, error)
	Delete(ctx context.Context, id int64) error
	Update(ctx context.Context, post *Post) error
	GetUserFeed(ctx context.Context, userID int64, fq PaginatedFeedQuery) ([]PostWithMetadata, error)
}

type UsersRepository interface {
	Create(ctx context.Context, tx *sql.Tx, user *User) error
	GetByID(ctx context.Context, id int64) (*User, error)
	CreateAndInvite(ctx context.Context, user *User, token string, invitationExp time.Duration) error
	Activate(ctx context.Context, token string) error
}

type CommentsRepository interface {
	GetByPostID(ctx context.Context, postID int64) ([]Comment, error)
	Create(ctx context.Context, comment *Comment) error
}

type FollowersRepository interface {
	Follow(ctx context.Context, userID, followerID int64) error
	Unfollow(ctx context.Context, userID, followerID int64) error
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
	Posts     PostsRepository
	Users     UsersRepository
	Comments  CommentsRepository
	Followers FollowersRepository
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
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &Repository{
		Posts:     &PostStore{db},
		Users:     &UserStore{db},
		Comments:  &CommentRepo{db},
		Followers: &FollowerRepo{db},
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

func withTx(db *sql.DB, ctx context.Context, fn func(tx *sql.Tx) error) (err error) {
	tx, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return err
	}
	if err := fn(tx); err != nil {
		return tx.Rollback()
	}
	return tx.Commit()
}
