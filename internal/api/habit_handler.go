package api

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/Fozzyack/habit-tracker/internal/models"
	"github.com/Fozzyack/habit-tracker/internal/services"
	"github.com/go-chi/chi"
	"github.com/rs/zerolog"
)

type HabitHandler struct {
	Logger       zerolog.Logger
	HabitService *services.HabitService
}

func NewHabitHandler(habitService *services.HabitService, logger zerolog.Logger) *HabitHandler {
	return &HabitHandler{
		Logger:       logger,
		HabitService: habitService,
	}
}

func (hh *HabitHandler) HandleCreateHabit(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	session := ctx.Value("session").(*models.Session)

	var habitReq *models.NewHabitRequest
	err := DecodeJSON(r, &habitReq)
	if err != nil {
		hh.Logger.Error().Err(err).Msg("HandleCreateHabit - Could not decode body")
		ErrorJSON(w, "Could not Create Habit (Bad Request)", http.StatusBadRequest)
		return
	}

	habit, err := hh.HabitService.CreateHabit(ctx, session.UserId, habitReq)
	if err != nil {
		hh.Logger.Error().Err(err).Msg("HandleCreateHabit - Could not create habit")
		ErrorJSON(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	SendJSON(w, habit)
}

func (hh *HabitHandler) HandleGetHabitById(w http.ResponseWriter, r *http.Request) {
	habitId := chi.URLParam(r, "habitId")

	habit, err := hh.HabitService.GetHabitById(habitId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ErrorJSON(w, "Habit Not Found", http.StatusNotFound)
			return
		}
		hh.Logger.Error().Err(err).Msg("HandleGetHabitById")
		ErrorJSON(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	SendJSON(w, habit)
}

func (hh *HabitHandler) HandleUpdateHabit(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	session := ctx.Value("session").(*models.Session)
	habitId := chi.URLParam(r, "habitId")

	var habitReq *models.UpdateHabitRequest
	err := DecodeJSON(r, &habitReq)
	if err != nil {
		hh.Logger.Error().Err(err).Msg("HandleUpdateHabit - Could not decode body")
		ErrorJSON(w, "Could not Update Habit (Bad Request)", http.StatusBadRequest)
		return
	}

	habit, err := hh.HabitService.UpdateHabit(ctx, habitId, session.UserId, habitReq)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ErrorJSON(w, "Habit Not Found", http.StatusNotFound)
			return
		}
		hh.Logger.Error().Err(err).Msg("HandleUpdateHabit - Could not update habit")
		ErrorJSON(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	SendJSON(w, habit)
}
