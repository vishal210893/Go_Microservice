package main

import (
	"Go-Microservice/internal/repo"
	"context"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

type userKey string

const userCtx userKey = "user"

type FollowerUser struct {
	UserId int64 `json:"userId"`
}

type FollowedUser struct {
	UserID     int64 `json:"userId"`
	FollowerId int64 `json:"followerId"`
}

func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromContext(r)
	if err := app.jsonResponse(w, http.StatusOK, user); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) followUserHandler(w http.ResponseWriter, r *http.Request) {
	followerUser := getUserFromContext(r)

	var payload FollowedUser
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	ctx := r.Context()
	if err := app.repo.Followers.Follow(ctx, followerUser.ID, payload.UserID); err != nil {
		switch err {
		case repo.ErrConflict:
			app.conflictResponse(w, r, err)
			return
		default:
			app.internalServerError(w, r, err)
			return
		}
	}

	if err := app.jsonResponse(w, http.StatusNoContent, nil); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) unfollowUserHandler(w http.ResponseWriter, r *http.Request) {
	unFollowingUser := getUserFromContext(r)

	var payload FollowedUser
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	ctx := r.Context()

	if err := app.repo.Followers.Unfollow(ctx, unFollowingUser.ID, payload.UserID); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusNoContent, nil); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) usersContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, err := strconv.ParseInt(chi.URLParam(r, "userID"), 10, 64)
		if err != nil {
			app.badRequestResponse(w, r, err)
			return
		}

		user, err := app.repo.Users.GetByID(r.Context(), userID)
		if err != nil {
			switch err {
			case repo.ErrPostNotFound:
				app.notFoundResponse(w, r, err)
				return
			default:
				app.internalServerError(w, r, err)
				return
			}
		}

		ctx := context.WithValue(r.Context(), userCtx, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getUserFromContext(r *http.Request) *repo.User {
	user, _ := r.Context().Value(userCtx).(*repo.User)
	return user
}
