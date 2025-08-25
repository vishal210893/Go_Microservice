package main

import (
	"Go-Microservice/internal/repo"
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type post string

const postCtx = "post"

type CreatePostPayload struct {
	Title   string   `json:"title" validate:"required,min=3,max=100"`
	Content string   `json:"content" validate:"required,max=100"`
	Tags    []string `json:"tags"`
	UserID  int64    `json:"userId" validate:"required"`
}

// CreatePost creates a new post with the provided content
//
//	@Summary		Create a new post
//	@Description	Create a new post with title, content, tags, and associate it with a user
//	@Description	All fields are validated before creation. Title must be between 3-100 characters.
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Param			post	body		CreatePostPayload	true	"Post creation payload"
//	@Success		201		{object}	repo.Post			"Post created successfully"
//	@Failure		400		{object}	map[string]string	"Invalid request payload or validation error"
//	@Failure		422		{object}	map[string]string	"Validation failed for post data"
//	@Failure		500		{object}	map[string]string	"Internal server error"
//	@Router			/v1/posts [post]
func (app *application) createPostHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreatePostPayload

	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	post := &repo.Post{
		Title:   payload.Title,
		Content: payload.Content,
		Tags:    payload.Tags,
		UserID:  payload.UserID,
	}

	ctx := r.Context()
	if err := app.repo.Posts.Create(ctx, post); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusCreated, post); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

// GetPost retrieves a specific post by ID with all comments
//
//	@Summary		Get post by ID
//	@Description	Retrieve detailed information about a specific post including all associated comments
//	@Description	Returns the complete post object with nested comments and user information
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Param			postID	path		int64				true	"Post ID to retrieve"
//	@Success		200		{object}	repo.Post			"Post details with comments retrieved successfully"
//	@Failure		400		{object}	map[string]string	"Invalid post ID format"
//	@Failure		404		{object}	map[string]string	"Post not found"
//	@Failure		500		{object}	map[string]string	"Internal server error"
//	@Router			/v1/posts/{postID} [get]
func (app *application) getPostHandler(w http.ResponseWriter, r *http.Request) {
	post := getPostFromCtx(r)

	comments, err := app.repo.Comments.GetByPostID(r.Context(), post.ID)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
	post.Comments = comments

	if err := app.jsonResponse(w, http.StatusOK, post); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

// DeletePost removes a post from the system
//
//	@Summary		Delete a post
//	@Description	Permanently delete a post and all associated data including comments
//	@Description	This action cannot be undone. Only the post owner can delete their posts.
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Param			postID	path	int64	true	"Post ID to delete"
//	@Success		204		"Post deleted successfully (no content returned)"
//	@Failure		400		{object}	map[string]string	"Invalid post ID format"
//	@Failure		404		{object}	map[string]string	"Post not found"
//	@Failure		500		{object}	map[string]string	"Internal server error"
//	@Router			/v1/posts/{postID} [delete]
func (app *application) deletePostHandler(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "postID")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	ctx := r.Context()

	if err := app.repo.Posts.Delete(ctx, id); err != nil {
		switch {
		case errors.Is(err, repo.ErrPostNotFound):
			app.notFoundResponse(w, r, err)
		default:

			app.internalServerError(w, r, err)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

type UpdatePostPayload struct {
	Title   *string  `json:"title" validate:"omitempty,min=3,max=100"`
	Content *string  `json:"content" validate:"omitempty,max=1000"`
	Tags    []string `json:"tags"`
}

// DeletePost removes a post from the system
//
//	@Summary		Delete a post
//	@Description	Permanently delete a post and all associated data including comments
//	@Description	This action cannot be undone. Only the post owner can delete their posts.
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Param			postID	path	int64	true	"Post ID to delete"
//	@Success		204		"Post deleted successfully (no content returned)"
//	@Failure		400		{object}	map[string]string	"Invalid post ID format"
//	@Failure		404		{object}	map[string]string	"Post not found"
//	@Failure		500		{object}	map[string]string	"Internal server error"
//	@Router			/v1/posts/{postID} [delete]
func (app *application) updatePostHandler(w http.ResponseWriter, r *http.Request) {
	post := getPostFromCtx(r)

	var payload UpdatePostPayload

	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := validate.Struct(&payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if payload.Title != nil {
		post.Title = *payload.Title
	}

	if payload.Content != nil {
		post.Content = *payload.Content
	}

	if payload.Tags != nil {
		post.Tags = payload.Tags
	}

	if err := app.repo.Posts.Update(r.Context(), post); err != nil {
		app.internalServerError(w, r, err)
		return
	}

}

func (app *application) postsContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idParam := chi.URLParam(r, "postID")
		id, err := strconv.ParseInt(idParam, 10, 64)
		if err != nil {
			app.internalServerError(w, r, err)
			return
		}

		ctx := r.Context()

		post, err := app.repo.Posts.GetByID(ctx, id)
		if err != nil {
			switch {
			case errors.Is(err, repo.ErrPostNotFound):
				app.notFoundResponse(w, r, err)
			default:
				app.internalServerError(w, r, err)
			}
			return
		}

		ctx = context.WithValue(ctx, postCtx, post)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getPostFromCtx(r *http.Request) *repo.Post {
	post, _ := r.Context().Value(postCtx).(*repo.Post)
	return post
}
