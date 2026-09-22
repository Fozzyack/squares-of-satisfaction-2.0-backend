package api

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/Fozzyack/habit-tracker/internal/models"
	"github.com/Fozzyack/habit-tracker/internal/services"
	"github.com/rs/zerolog"
)

type UserHandler struct {
	Logger      zerolog.Logger
	AuthService *services.AuthService
	UserService *services.UserService
}

func NewUserHandler(authService *services.AuthService, userService *services.UserService, logger zerolog.Logger) *UserHandler {
	return &UserHandler{
		Logger:      logger,
		AuthService: authService,
		UserService: userService,
	}
}

func (uh *UserHandler) HandleGetUserBySession(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	session := ctx.Value("session").(*models.Session)
	user, err := uh.AuthService.UserStore.GetUserById(ctx, session.UserId)
	if err != nil {
		uh.Logger.Error().Err(err).Msg("HandleGetUserBySession")
		ErrorJSON(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	SendJSON(w, map[string]string{"name": user.Name, "email": user.Email})
}

func (uh *UserHandler) HandleUpdateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	session := ctx.Value("session").(*models.Session)

	var userReq *models.UpdateUserRequest
	err := DecodeJSON(r, &userReq)
	if err != nil || userReq == nil || userReq.CurrentPassword == "" {
		if err != nil {
			uh.Logger.Error().Err(err).Msg("HandleUpdateUser - Could not decode body")
		}
		ErrorJSON(w, "Could not Update User (Bad Request)", http.StatusBadRequest)
		return
	}

	user, err := uh.AuthService.UpdateUser(ctx, session.UserId, userReq)
	if err != nil {
		if errors.Is(err, services.ErrInvalidCredentials) {
			ErrorJSON(w, "Incorrect Current Password", http.StatusUnauthorized)
			return
		}
		if errors.Is(err, sql.ErrNoRows) {
			ErrorJSON(w, "User Not Found", http.StatusNotFound)
			return
		}
		uh.Logger.Error().Err(err).Msg("HandleUpdateUser - Could not update user")
		ErrorJSON(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	SendJSON(w, user)
}

func (uh *UserHandler) HandleDeleteUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	session := ctx.Value("session").(*models.Session)

	err := uh.UserService.DeleteUser(ctx, session.UserId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ErrorJSON(w, "User Not Found", http.StatusNotFound)
			return
		}
		uh.Logger.Error().Err(err).Msg("HandleDeleteUser - Could not delete user")
		ErrorJSON(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	cookie, err := DeleteCookie()
	if err != nil {
		uh.Logger.Error().Err(err).Msg("HandleDeleteUser - cookie name")
		ErrorJSON(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, cookie)

	SendJSON(w, map[string]string{"msg": "success"})
}

func (uh *UserHandler) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var userReq *models.NewUserRequest
	err := DecodeJSON(r, &userReq)
	if err != nil {
		uh.Logger.Error().Err(err).Msg("HandleCreateUser - Could not decode body")
		ErrorJSON(w, "Could not Create User (Bad Request)", http.StatusBadRequest)
		return
	}

	_, session, err := uh.AuthService.CreateNewUser(ctx, userReq)
	if err != nil {
		uh.Logger.Error().Err(err).Msg("HandleCreateUser - Could not create user / session")
		ErrorJSON(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	cookie, err := CreateCookie(session.Token, session.ExpiresAt)
	if err != nil {
		uh.Logger.Error().Err(err).Msg("HandleCreateUser - cookie name")
		ErrorJSON(w, "Internal Server Error", http.StatusInternalServerError)
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
		ErrorJSON(w, "Could not Login User (Bad Request)", http.StatusBadRequest)
		return
	}

	session, err := uh.AuthService.LoginUser(ctx, loginReq)
	if err != nil {
		if errors.Is(err, services.ErrInvalidCredentials) {
			ErrorJSON(w, "Incorrect Email or Password", http.StatusUnauthorized)
			return
		}
		uh.Logger.Error().Err(err).Msg("HandleLoginUser - Could not login user")
		ErrorJSON(w, "Could not Login User (Bad Request)", http.StatusBadRequest)
		return
	}

	cookie, err := CreateCookie(session.Token, session.ExpiresAt)
	if err != nil {
		uh.Logger.Error().Err(err).Msg("HandleCreateUser - cookie name")
		ErrorJSON(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, cookie)

	SendJSON(w, map[string]string{"msg": "success"})
}
