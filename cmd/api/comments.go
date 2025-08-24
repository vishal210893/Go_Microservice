package main

import (
	"Go-Microservice/internal/repo"
	"net/http"
)

type CreateCommentPayload struct {
	Content string `json:"content" validate:"required"`
	UserID  int64  `json:"userId" validate:"required"`
	Likes   int    `json:"likes"`
}

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
