// Package repo provides data access layer implementations for application entities.
// It contains repository patterns for database operations with PostgreSQL.
package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/lib/pq"
)

// Query timeout constants
const (
	// DefaultQueryTimeout is the default timeout for database queries
	DefaultQueryTimeout = 15 * time.Second

	// DefaultTransactionTimeout is the default timeout for database transactions
	DefaultTransactionTimeout = 30 * time.Second
)

// Common repository errors
var (
	ErrPostNotFound    = errors.New("post not found")
	ErrPostExists      = errors.New("post already exists")
	ErrInvalidPostData = errors.New("invalid post data")
)

// Post represents a blog post entity with metadata and content.
// It maps to the posts table in the database and includes audit fields.
type Post struct {
	ID        int64     `json:"id" db:"id"`
	Content   string    `json:"content" db:"content"`
	Title     string    `json:"title" db:"title"`
	UserID    int64     `json:"user_id" db:"user_id"`
	Tags      []string  `json:"tags" db:"tags"`
	CreatedAt string    `json:"created_at" db:"created_at"`
	UpdatedAt string    `json:"updated_at" db:"updated_at"`
	Comments  []Comment `json:"comments" db:"comments"`
}

// Validate checks if the post data is valid for database operations
func (p *Post) Validate() error {
	if p.Title == "" {
		return fmt.Errorf("%w: title is required", ErrInvalidPostData)
	}
	if p.Content == "" {
		return fmt.Errorf("%w: content is required", ErrInvalidPostData)
	}
	if p.UserID <= 0 {
		return fmt.Errorf("%w: valid user_id is required", ErrInvalidPostData)
	}
	return nil
}

// PostStore provides database operations for Post entities.
// It implements the repository pattern for posts table operations.
type PostStore struct {
	db *sql.DB
}

// NewPostStore creates a new PostStore with the provided database connection.
// It uses default timeout values for database operations.
func NewPostStore(db *sql.DB) *PostStore {
	return &PostStore{
		db: db,
	}
}

// Create inserts a new post into the database.
// It validates the post data, executes the insert query, and returns the created post
// with populated ID and timestamps.
//
// The function will:
//   - Validate post data before insertion
//   - Execute INSERT query with RETURNING clause
//   - Handle database constraint violations
//   - Return error if validation or insertion fails
//
// Example:
//
//	post := &Post{
//		Title:   "My Post",
//		Content: "Post content here",
//		UserID:  123,
//		Tags:    []string{"go", "database"},
//	}
//	err := store.Create(ctx, post)
//	if err != nil {
//		log.Printf("Failed to create post: %v", err)
//	}
func (postStore *PostStore) Create(ctx context.Context, post *Post) error {
	if post == nil {
		return fmt.Errorf("%w: post cannot be nil", ErrInvalidPostData)
	}

	if err := post.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	query := `
		INSERT INTO posts (content, title, user_id, tags)
		VALUES ($1, $2, $3, $4) 
		RETURNING id, created_at, updated_at
	`

	// Set query timeout
	//ctx, cancel := context.WithTimeout(ctx, postStore.queryTimeout)
	//defer cancel()

	err := postStore.db.QueryRowContext(
		ctx,
		query,
		post.Content,
		post.Title,
		post.UserID,
		pq.Array(post.Tags),
	).Scan(
		&post.ID,
		&post.CreatedAt,
		&post.UpdatedAt,
	)

	if err != nil {
		// Handle PostgreSQL specific errors
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code {
			case "23505": // unique_violation
				return fmt.Errorf("%w: %v", ErrPostExists, pqErr.Detail)
			case "23503": // foreign_key_violation
				return fmt.Errorf("foreign key constraint violation: %w", err)
			case "23514": // check_violation
				return fmt.Errorf("check constraint violation: %w", err)
			}
		}

		// Handle context errors
		if errors.Is(err, context.DeadlineExceeded) {
			return fmt.Errorf("create operation timed out: %w", err)
		}

		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("no rows returned from insert: %w", err)
		}

		return fmt.Errorf("failed to create post: %w", err)
	}

	return nil
}

// GetByID retrieves a post by its ID from the database.
// Returns ErrPostNotFound if the post doesn't exist.
func (s *PostStore) GetByID(ctx context.Context, id int64) (*Post, error) {
	if id <= 0 {
		return nil, fmt.Errorf("%w: invalid post ID", ErrInvalidPostData)
	}

	query := `
		SELECT id, content, title, user_id, tags, created_at, updated_at
		FROM posts
		WHERE id = $1
	`

	//ctx, cancel := context.WithTimeout(ctx, s.queryTimeout)
	//defer cancel()

	post := Post{}
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&post.ID,
		&post.Content,
		&post.Title,
		&post.UserID,
		pq.Array(&post.Tags),
		&post.CreatedAt,
		&post.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%w: post with ID %d", ErrPostNotFound, id)
		}
		return nil, fmt.Errorf("failed to get post: %w", err)
	}

	return &post, nil
}

//
//// Update modifies an existing post using optimistic locking.
//// The version field is used to prevent concurrent modification conflicts.
//func (s *PostStore) Update(ctx context.Context, post *Post) error {
//	if post == nil {
//		return fmt.Errorf("%w: post cannot be nil", ErrInvalidPostData)
//	}
//
//	if err := post.Validate(); err != nil {
//		return fmt.Errorf("validation failed: %w", err)
//	}
//
//	if post.ID <= 0 {
//		return fmt.Errorf("%w: invalid post ID for update", ErrInvalidPostData)
//	}
//
//	query := `
//		UPDATE posts
//		SET content = $1, title = $2, tags = $3, updated_at = CURRENT_TIMESTAMP, version = version + 1
//		WHERE id = $4 AND version = $5
//		RETURNING updated_at, version
//	`
//
//	ctx, cancel := context.WithTimeout(ctx, s.queryTimeout)
//	defer cancel()
//
//	err := s.db.QueryRowContext(
//		ctx,
//		query,
//		post.Content,
//		post.Title,
//		pq.Array(post.Tags),
//		post.ID,
//		post.Version,
//	).Scan(&post.UpdatedAt, &post.Version)
//
//	if err != nil {
//		if errors.Is(err, sql.ErrNoRows) {
//			return fmt.Errorf("%w: post may not exist or version conflict", ErrPostNotFound)
//		}
//		return fmt.Errorf("failed to update post: %w", err)
//	}
//
//	return nil
//}
//
//// Delete removes a post from the database by ID.
//// Returns ErrPostNotFound if the post doesn't exist.
//func (s *PostStore) Delete(ctx context.Context, id int64) error {
//	if id <= 0 {
//		return fmt.Errorf("%w: invalid post ID", ErrInvalidPostData)
//	}
//
//	query := `DELETE FROM posts WHERE id = $1`
//
//	ctx, cancel := context.WithTimeout(ctx, s.queryTimeout)
//	defer cancel()
//
//	result, err := s.db.ExecContext(ctx, query, id)
//	if err != nil {
//		return fmt.Errorf("failed to delete post: %w", err)
//	}
//
//	rowsAffected, err := result.RowsAffected()
//	if err != nil {
//		return fmt.Errorf("failed to get rows affected: %w", err)
//	}
//
//	if rowsAffected == 0 {
//		return fmt.Errorf("%w: post with ID %d", ErrPostNotFound, id)
//	}
//
//	return nil
//}
