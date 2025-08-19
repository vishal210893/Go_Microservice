package main

import (
	"Go-Microservice/internal/repo"
	"errors"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

type CreatePostPayload struct {
	Title   string   `json:"title"`
	Content string   `json:"content"`
	Tags    []string `json:"tags"`
}

func (app *application) createPostHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreatePostPayload

	err := readJSON(w, r, &payload)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
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
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := writeJSON(w, http.StatusCreated, post); err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

}

func (app *application) getPostHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	postID := chi.URLParam(r, "postID")
	postId, err2 := strconv.ParseInt(postID, 10, 64)
	if err2 != nil {
		writeJSONError(w, http.StatusBadRequest, err2.Error())
		return
	}
	post, err := app.repo.Posts.GetByID(ctx, postId)
	if err != nil {
		switch {
		case errors.Is(err, repo.ErrPostNotFound):
			writeJSONError(w, http.StatusNotFound, err.Error())
		default:
			writeJSONError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	if err := writeJSON(w, http.StatusOK, post); err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
}
