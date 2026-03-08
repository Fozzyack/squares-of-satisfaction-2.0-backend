package routes

import (
	"context"
	"net/http"
	"time"

	"github.com/Fozzyack/habit-tracker/internal/api"
	"github.com/Fozzyack/habit-tracker/internal/app"
	"github.com/Fozzyack/habit-tracker/internal/env"
)

func CheckSession(app *app.Application) func(next http.Handler) http.Handler {
	app.Logger.Info().Msg("Running Session Check Middleware")
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookieName, err := env.GetCookieName()
			if err != nil {
				app.Logger.Error().Err(err).Msg("Check session: .env error")
				api.ErrorJSON(w, "Error: Internal Server Error", http.StatusInternalServerError)
				return
			}

			cookie, err := r.Cookie(cookieName)
			if err != nil {
				app.Logger.Error().Err(err).Msg("Check Session: Could not get cookie")
				api.ErrorJSON(w, "Error: Unauthorized", http.StatusUnauthorized)
				return
			}

			session, err := app.SessionStore.GetSessionByToken(cookie.Value)
			if err != nil {
				app.Logger.Error().Err(err).Msg("Check Session: Session not found")
				api.ErrorJSON(w, "Error: Unauthorized", http.StatusUnauthorized)
				return
			}

			if session.ExpiresAt.UTC().Before(time.Now().UTC()) {
				app.Logger.Error().Msg("Check Session: Session Expired")
				api.ErrorJSON(w, "Error: Unauthorized", http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), "session", session)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
