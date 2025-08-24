// Package repo provides data access layer implementations for application entities.
// It contains repository patterns for database operations with PostgreSQL.
package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/lib/pq"
)

var (
	ErrPostNotFound    = errors.New("post not found")
	ErrPostExists      = errors.New("post already exists")
	ErrInvalidPostData = errors.New("invalid post data")
)

type Post struct {
	ID        int64     `json:"id" db:"id"`
	Content   string    `json:"content" db:"content"`
	Title     string    `json:"title" db:"title"`
	UserID    int64     `json:"user_id" db:"user_id"`
	Tags      []string  `json:"tags" db:"tags"`
	CreatedAt string    `json:"created_at" db:"created_at"`
	UpdatedAt string    `json:"updated_at" db:"updated_at"`
	Version   int32     `json:"-" db:"version"`
	Comments  []Comment `json:"comments" db:"comments"`
	User      User      `json:"user" db:"user"`
}

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

type PostStore struct {
	db *sql.DB
}

type PostWithMetadata struct {
	Post
	CommentsCount int `json:"comments_count"`
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

	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()

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

func (postStore *PostStore) GetByID(ctx context.Context, id int64) (*Post, error) {
	if id <= 0 {
		return nil, fmt.Errorf("%w: invalid post ID", ErrInvalidPostData)
	}

	query := `
		SELECT id, content, title, user_id, tags, created_at, updated_at, version
		FROM posts
		WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()

	post := Post{}
	err := postStore.db.QueryRowContext(ctx, query, id).Scan(
		&post.ID,
		&post.Content,
		&post.Title,
		&post.UserID,
		pq.Array(&post.Tags),
		&post.CreatedAt,
		&post.UpdatedAt,
		&post.Version,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%w: post with ID %d", ErrPostNotFound, id)
		}
		return nil, fmt.Errorf("failed to get post: %w", err)
	}

	return &post, nil
}

func (postStore *PostStore) Update(ctx context.Context, post *Post) error {
	if post == nil {
		return fmt.Errorf("%w: post cannot be nil", ErrInvalidPostData)
	}
	if err := post.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}
	if post.ID <= 0 {
		return fmt.Errorf("%w: invalid post ID for update", ErrInvalidPostData)
	}

	query := `
		UPDATE posts
		SET title = $1, content = $2, tags = $3, updated_at = NOW(), version = version + 1
		WHERE id = $4 AND version = $5
		RETURNING version
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()

	err := postStore.db.QueryRowContext(
		ctx,
		query,
		post.Title,
		post.Content,
		pq.Array(post.Tags),
		post.ID,
		post.Version,
	).Scan(&post.Version)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("%w: post may not exist or version conflict", ErrPostNotFound)
		}
		return fmt.Errorf("failed to update post: %w", err)
	}
	return nil
}

func (postStore *PostStore) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return fmt.Errorf("%w: invalid post ID", ErrInvalidPostData)
	}

	query := `DELETE FROM posts WHERE id = $1`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()

	result, err := postStore.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete post: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("%w: post with ID %d", ErrPostNotFound, id)
	}

	return nil
}

func (postStore *PostStore) GetUserFeed(ctx context.Context, userID int64) ([]PostWithMetadata, error) {
	query := `
		SELECT 
			p.id, p.user_id, p.title, p.content, p.created_at, p.version, p.tags,
			u.username,
			COUNT(c.id) AS comments_count
		FROM posts p
		LEFT JOIN comments c ON c.post_id = p.id
		LEFT JOIN users u ON p.user_id = u.id
		JOIN followers f ON f.follower_id = p.user_id OR p.user_id = $1
		WHERE f.user_id = $1 OR p.user_id = $1 
		GROUP BY p.id, u.username,p.created_at
		ORDER BY p.created_at DESC
		`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()

	rows, err := postStore.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var feed []PostWithMetadata
	for rows.Next() {
		var p PostWithMetadata
		err := rows.Scan(
			&p.ID,
			&p.UserID,
			&p.Title,
			&p.Content,
			&p.CreatedAt,
			&p.Version,
			pq.Array(&p.Tags),
			&p.User.Username,
			&p.CommentsCount,
		)
		if err != nil {
			return nil, err
		}
		feed = append(feed, p)
	}
	return feed, nil
}
