package main

import (
	"Go-Microservice/internal/repo"
	"net/http"
)

// CreateCommentPayload represents the request payload for creating a comment
//
//	@Description	Request payload for creating a new comment
type CreateCommentPayload struct {
	// Comment content
	//	@example	"This is a great post!"
	Content string `json:"content" validate:"required,min=1,max=500" example:"This is a great post!"`

	// ID of the user creating the comment
	//	@example	1
	UserID int64 `json:"user_id" validate:"required" example:"1"`

	// No. of likes
	//	@example	5
	Likes int `json:"likes"`
}

// CreateComment adds a new comment to a specific post
//
//	@Summary		Create a comment on a post
//	@Description	Add a new comment to an existing post with the provided content and user information
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Param			postID	path		int64					true	"Post ID to comment on"
//	@Param			comment	body		CreateCommentPayload	true	"Comment creation payload"
//	@Success		201		{object}	repo.Comment			"Comment created successfully"
//	@Failure		400		{object}	map[string]string		"Invalid request payload or post ID"
//	@Failure		404		{object}	map[string]string		"Post not found"
//	@Failure		422		{object}	map[string]string		"Validation failed for comment data"
//	@Failure		500		{object}	map[string]string		"Internal server error"
//	@Router			/v1/posts/{postID}/comments [post]
func (app *application) createCommentHandler(w http.ResponseWriter, r *http.Request) {
	post := getPostFromCtx(r)

	var payload CreateCommentPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	comment := &repo.Comment{
		PostID:  post.ID,
		UserID:  payload.UserID,
		Content: payload.Content,
	}

	if err := app.repo.Comments.Create(r.Context(), comment); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusCreated, comment); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
