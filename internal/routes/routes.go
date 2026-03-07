package routes

import (
	"github.com/Fozzyack/habit-tracker/internal/app"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func SetupRoutes(app *app.Application) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)

	return r;
}
