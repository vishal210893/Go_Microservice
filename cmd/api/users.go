package main

import (
	"Go-Microservice/internal/repo"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type userKey string

const userCtx userKey = "user"

// GetUser retrieves a specific user by their ID
//
//	@Summary		Get user by ID
//	@Description	Retrieve detailed information about a specific user including their profile data
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			userID	path		int64				true	"User ID to retrieve"
//	@Success		200		{object}	repo.User			"User details retrieved successfully"
//	@Failure		400		{object}	map[string]string	"Invalid user ID format"
//	@Failure		404		{object}	map[string]string	"User not found"
//	@Failure		500		{object}	map[string]string	"Internal server error"
//	@Router			/v1/users/{userID} [get]
func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseInt(chi.URLParam(r, "userID"), 10, 64)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user, err := app.repo.Users.GetByID(r.Context(), userID)
	if err != nil {
		switch err {
		case repo.ErrNotFound:
			app.notFoundResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}
	if err := app.jsonResponse(w, http.StatusOK, user); err != nil {
		app.internalServerError(w, r, err)
	}
}

// FollowUser allows the current user to follow another user
//
//	@Summary		Follow a user
//	@Description	Create a following relationship between the authenticated user and the target user
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			userID	path		int64				true	"ID of the user to follow"
//	@Success		200		{object}	map[string]string	"Successfully followed user"
//	@Success		201		{object}	map[string]string	"Following relationship created"
//	@Failure		400		{object}	map[string]string	"Invalid user ID or cannot follow yourself"
//	@Failure		404		{object}	map[string]string	"User not found"
//	@Failure		409		{object}	map[string]string	"Already following this user"
//	@Failure		500		{object}	map[string]string	"Internal server error"
//	@Security		BasicAuth
//	@Router			/v1/users/{userID}/follow [put]
func (app *application) followUserHandler(w http.ResponseWriter, r *http.Request) {
	followingUser := getUserFromContext(r)
	followedID, err := strconv.ParseInt(chi.URLParam(r, "userID"), 10, 64)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	ctx := r.Context()
	if err := app.repo.Followers.Follow(ctx, followingUser.ID, followedID); err != nil {
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

// UnfollowUser allows the current user to unfollow another user
//
//	@Summary		Unfollow a user
//	@Description	Remove the following relationship between the authenticated user and the target user
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			userID	path		int64				true	"ID of the user to unfollow"
//	@Success		200		{object}	map[string]string	"Successfully unfollowed user"
//	@Success		204		{object}	map[string]string	"Following relationship removed"
//	@Failure		400		{object}	map[string]string	"Invalid user ID or cannot unfollow yourself"
//	@Failure		404		{object}	map[string]string	"User not found or not following"
//	@Failure		500		{object}	map[string]string	"Internal server error"
//	@Security		BasicAuth
//	@Router			/v1/users/{userID}/unfollow [put]
func (app *application) unfollowUserHandler(w http.ResponseWriter, r *http.Request) {
	unFollowingUser := getUserFromContext(r)

	unFollowedID, err := strconv.ParseInt(chi.URLParam(r, "userID"), 10, 64)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	ctx := r.Context()

	if err := app.repo.Followers.Unfollow(ctx, unFollowingUser.ID, unFollowedID); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusNoContent, nil); err != nil {
		app.internalServerError(w, r, err)
	}
}

// ActivateUser godoc
//
//	@Summary		Activates/Register a user
//	@Description	Activates/Register a user by invitation token
//	@Tags			users
//	@Produce		json
//	@Param			token	path		string	true	"Invitation token"
//	@Success		204		{string}	string	"User activated"
//	@Failure		404		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/users/activate/{token} [put]
func (app *application) activateUserHandler(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")

	err := app.repo.Users.Activate(r.Context(), token)
	if err != nil {
		switch err {
		case repo.ErrNotFound:
			app.notFoundResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	if err := app.jsonResponse(w, http.StatusNoContent, ""); err != nil {
		app.internalServerError(w, r, err)
	}
}

func getUserFromContext(r *http.Request) *repo.User {
	user, _ := r.Context().Value(userCtx).(*repo.User)
	return user
}
