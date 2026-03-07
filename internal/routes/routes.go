package routes

import (
	"encoding/json"
	"net/http"

	"github.com/Fozzyack/habit-tracker/internal/app"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"msg": "Server is Healthy!"})
	return
}

func SetupRoutes(app *app.Application) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)

	r.Get("/health", healthCheck)
	r.Post("/user", app.UserHandler.HandleCreateUser)

	return r
}
