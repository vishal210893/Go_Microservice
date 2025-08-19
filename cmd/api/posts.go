package main

import (
	"Go-Microservice/internal/repo"
	"errors"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

type CreatePostPayload struct {
	Title   string   `json:"title" validate:"required,min=3,max=100"`
	Content string   `json:"content" validate:"required,max=100"`
	Tags    []string `json:"tags"`
}

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
		UserID:  1,
	}

	ctx := r.Context()
	if err := app.repo.Posts.Create(ctx, post); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := writeJSON(w, http.StatusCreated, post); err != nil {
		app.internalServerError(w, r, err)
		return
	}

}

func (app *application) getPostHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	postID := chi.URLParam(r, "postID")
	postId, err := strconv.ParseInt(postID, 10, 64)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	post, err := app.repo.Posts.GetByID(ctx, postId)
	if err != nil {
		switch {
		case errors.Is(err, repo.ErrPostNotFound):
			app.notFoundResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	comments, err := app.repo.Comments.GetByPostID(ctx, postId)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	post.Comments = comments
	
	if err := writeJSON(w, http.StatusOK, post); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
