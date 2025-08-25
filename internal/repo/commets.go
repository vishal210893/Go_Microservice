package repo

import (
	"context"
	"database/sql"
)

// Comment represents a comment on a post
//
//	@Description	User comment on a specific post
type Comment struct {
	// Comment's unique identifier
	//	@example	1
	ID int64 `json:"id" example:"1"`

	// Comment content
	//	@example	"Great post! Thanks for sharing."
	Content string `json:"content" example:"Great post! Thanks for sharing."`

	// ID of the post this comment belongs to
	//	@example	1
	PostID int64 `json:"post_id" example:"1"`

	// ID of the user who created the comment
	//	@example	2
	UserID int64 `json:"user_id" example:"2"`

	// Timestamp when the comment was created
	//	@example	2024-01-15T11:30:00Z
	CreatedAt string `json:"created_at" example:"2024-01-15T11:30:00Z"`

	User User `json:"user"`
}

type CommentRepo struct {
	db *sql.DB
}

func (commentRepo *CommentRepo) GetByPostID(ctx context.Context, postID int64) ([]Comment, error) {
	query := `
		SELECT c.id, c.post_id, c.user_id, c.content, c.created_at, users.username, users.id  FROM comments c
		JOIN users on users.id = c.user_id
		WHERE c.post_id = $1
		ORDER BY c.created_at DESC;
	`

	rows, err := commentRepo.db.QueryContext(ctx, query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	comments := []Comment{}
	for rows.Next() {
		var c Comment
		c.User = User{}
		err := rows.Scan(&c.ID, &c.PostID, &c.UserID, &c.Content, &c.CreatedAt, &c.User.Username, &c.User.ID)
		if err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}

	return comments, nil
}

func (commentRepo *CommentRepo) Create(ctx context.Context, comment *Comment) error {
	query := `
		INSERT INTO comments (post_id, user_id, content)
		VALUES ($1, $2, $3)
		RETURNING id, created_at
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeout)
	defer cancel()

	err := commentRepo.db.QueryRowContext(
		ctx,
		query,
		comment.PostID,
		comment.UserID,
		comment.Content,
	).Scan(
		&comment.ID,
		&comment.CreatedAt,
	)
	if err != nil {
		return err
	}
	return nil
}
