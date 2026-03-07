package api

import (
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
		ErrorJSON(w, "Error: Could not Create User", http.StatusBadRequest)
		return
	}

	SendJSON(w, map[string]string{"token": session.Token})
}
