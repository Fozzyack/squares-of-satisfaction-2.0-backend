package api

import (
	"errors"
	"net/http"

	"github.com/Fozzyack/habit-tracker/internal/models"
	"github.com/Fozzyack/habit-tracker/internal/services"
	"github.com/rs/zerolog"
)

type UserHandler struct {
	Logger      zerolog.Logger
	AuthService *services.AuthService
}

func NewUserHandler(authService *services.AuthService, logger zerolog.Logger) *UserHandler {
	return &UserHandler{
		Logger:      logger,
		AuthService: authService,
	}
}

func (uh *UserHandler) HandleGetUserBySession(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	session := ctx.Value("session").(*models.Session)
	user, err := uh.AuthService.UserStore.GetUserById(session.UserId)
	if err != nil {
		uh.Logger.Error().Err(err).Msg("HandleGetUserBySession")
		ErrorJSON(w, "Error: Internal Server Error", http.StatusInternalServerError)
		return
	}

	SendJSON(w, map[string]string{"name": user.Name, "email": user.Email})
}

func (uh *UserHandler) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var userReq *models.NewUserRequest
	err := DecodeJSON(r, &userReq)
	if err != nil {
		uh.Logger.Error().Err(err).Msg("HandleCreateUser - Could not decode body")
		ErrorJSON(w, "Error: Could not Create User", http.StatusBadRequest)
		return
	}

	_, session, err := uh.AuthService.CreateNewUser(ctx, userReq)
	if err != nil {
		uh.Logger.Error().Err(err).Msg("HandleCreateUser - Could not create user / session")
		ErrorJSON(w, "Error: Could not Create User", http.StatusInternalServerError)
		return
	}

	cookie, err := CreateCookie(session.Token, session.ExpiresAt)
	if err != nil {
		uh.Logger.Error().Err(err).Msg("HandleCreateUser - cookie name")
		ErrorJSON(w, "Error: Internal Server Error", http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, cookie)
	SendJSON(w, map[string]string{"msg": "success"})
}

func (uh *UserHandler) HandleLoginUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var loginReq *models.LoginUserRequest
	err := DecodeJSON(r, &loginReq)
	if err != nil {
		uh.Logger.Error().Err(err).Msg("HandleLoginUser - Could not decode body")
		ErrorJSON(w, "Error: Could not Login User", http.StatusBadRequest)
		return
	}

	session, err := uh.AuthService.LoginUser(ctx, loginReq)
	if err != nil {
		if errors.Is(err, services.ErrInvalidCredentials) {
			ErrorJSON(w, "Error: Invalid credentials", http.StatusUnauthorized)
			return
		}
		uh.Logger.Error().Err(err).Msg("HandleLoginUser - Could not login user")
		ErrorJSON(w, "Error: Could not Login User", http.StatusBadRequest)
		return
	}

	cookie, err := CreateCookie(session.Token, session.ExpiresAt)
	if err != nil {
		uh.Logger.Error().Err(err).Msg("HandleCreateUser - cookie name")
		ErrorJSON(w, "Error: Internal Server Error", http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, cookie)

	SendJSON(w, map[string]string{"msg": "success"})
}
