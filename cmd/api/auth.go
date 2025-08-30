package main

import (
	"Go-Microservice/internal/repo"
	"crypto/sha256"
	"encoding/hex"
	"github.com/google/uuid"
	"net/http"
)

type RegisterUserPayload struct {
	Username string `json:"username" validate:"required,max=100"`
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=3,max=72"`
}

type UserWithToken struct {
	*repo.User
	Token string `json:"token"`
}

// registerUserHandler godoc
//
//	@Summary		Registers a user
//	@Description	Registers a user
//	@Tags			authentication
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		RegisterUserPayload	true	"User credentials"
//	@Success		201		{object}	UserWithToken		"User registered"
//	@Failure		400		{object}	error
//	@Failure		500		{object}	error
//	@Router			/authentication/user [post]
func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var payload RegisterUserPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := &repo.User{
		Username: payload.Username,
		Email:    payload.Email,
	}

	// hash the user password
	if err := user.Password.Set(payload.Password); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	ctx := r.Context()

	plainToken := uuid.New().String()

	// hash the token for storage but keep the plain token for email
	hash := sha256.Sum256([]byte(plainToken))
	hashToken := hex.EncodeToString(hash[:])

	err := app.repo.Users.CreateAndInvite(ctx, user, hashToken, app.config.invitationExpTime)
	if err != nil {
		switch err {
		case repo.ErrDuplicateEmail:
			app.badRequestResponse(w, r, err)
		case repo.ErrDuplicateUsername:
			app.badRequestResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	userWithToken := UserWithToken{
		User:  user,
		Token: plainToken,
	}
	//activationURL := fmt.Sprintf("%s/confirm/%s", app.config.frontendURL, plainToken)
	//
	//isProdEnv := app.config.env == "production"
	//vars := struct {
	//	Username      string
	//	ActivationURL string
	//}{
	//	Username:      user.Username,
	//	ActivationURL: activationURL,
	//}
	//
	//// send mail
	//status, err := app.mailer.Send(mailer.UserWelcomeTemplate, user.Username, user.Email, vars, !isProdEnv)
	//if err != nil {
	//	app.logger.Errorw("error sending welcome email", "error", err)
	//
	//	// rollback user creation if email fails (SAGA pattern)
	//	if err := app.repo.Users.Delete(ctx, user.ID); err != nil {
	//		app.logger.Errorw("error deleting user", "error", err)
	//	}
	//
	//	app.internalServerError(w, r, err)
	//	return
	//}

	//app.logger.Info("Email sent", "status code", status)

	if err := app.jsonResponse(w, http.StatusCreated, userWithToken); err != nil {
		app.internalServerError(w, r, err)
	}
}
