package main

import (
	"Go-Microservice/internal/repo"
	"net/http"
)

// GetUserFeed retrieves the personalized feed for the authenticated user
//
//	@Summary		Get user's personalized feed
//	@Description	Retrieve a chronological feed of posts from users that the authenticated user follows
//	@Description	Returns posts ordered by creation date (newest first) with pagination support
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			limit	query		int						false	"Number of posts to retrieve (default: 20, max: 100)"
//	@Param			offset	query		int						false	"Number of posts to skip for pagination (default: 0)"
//	@Param			since	query		string					false	"ISO 8601 timestamp to get posts created after this time"
//	@Success		200		{array}		repo.Post				"User feed retrieved successfully"
//	@Success		200		{object}	map[string]interface{}	"Feed with pagination metadata"
//	@Failure		400		{object}	map[string]string		"Invalid query parameters"
//	@Failure		401		{object}	map[string]string		"Authentication required"
//	@Failure		500		{object}	map[string]string		"Internal server error"
//	@Security		BasicAuth
//	@Router			/v1/users/feed [get]
func (app *application) getUserFeedHandler(w http.ResponseWriter, r *http.Request) {
	fq := repo.PaginatedFeedQuery{
		Limit:  20,
		Offset: 0,
		Sort:   "desc",
		Tags:   []string{},
		Search: "",
	}

	fq, err := fq.Parse(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := validate.Struct(fq); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	ctx := r.Context()
	//user := getUserFromContext(r)

	feed, err := app.repo.Posts.GetUserFeed(ctx, int64(156), fq)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, feed); err != nil {
		app.internalServerError(w, r, err)
	}
}
