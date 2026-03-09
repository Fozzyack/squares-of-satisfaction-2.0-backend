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

func (hh *HabitHandler) HandleRecordHabit(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	session := ctx.Value("session").(*models.Session)
	habitId := chi.URLParam(r, "habitId")
	if habitId == "" {
		ErrorJSON(w, "Missing habit id", http.StatusBadRequest)
		return
	}

	var habitDailyReq *models.RecordHabitRequest
	err := DecodeJSON(r, &habitDailyReq)
	if err != nil {
		hh.Logger.Error().Err(err).Msg("HandleRecordHabit - Could not decode body")
		ErrorJSON(w, "Could not Record Habit (Bad Request)", http.StatusBadRequest)
		return
	}

	if habitDailyReq.HabitId != "" && habitDailyReq.HabitId != habitId {
		ErrorJSON(w, "habit_id in body must match URL parameter", http.StatusBadRequest)
		return
	}
	habitDailyReq.HabitId = habitId

	habitDailyTotal, err := hh.HabitService.RecordHabit(ctx, session.UserId, habitDailyReq)
	if err != nil {
		if errors.Is(err, services.ErrInvalidRecordHabitDate) {
			ErrorJSON(w, "Invalid date format. Use YYYY-MM-DD", http.StatusBadRequest)
			return
		}
		hh.Logger.Error().Err(err).Msg("HandleRecordHabit - Could not update / create total")
		ErrorJSON(w, "Error: Internal Server Error", http.StatusInternalServerError)
		return
	}

	SendJSON(w, habitDailyTotal)
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

func (hh *HabitHandler) HandleGetHabitsBySession(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	session := ctx.Value("session").(*models.Session)

	habits, err := hh.HabitService.GetHabitsByUserId(ctx, session.UserId)
	if err != nil {
		hh.Logger.Error().Err(err).Msg("HandleGetHabitsBySession")
		ErrorJSON(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	SendJSON(w, habits)
}

func (hh *HabitHandler) HandleGetHabitById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	session := ctx.Value("session").(*models.Session)
	habitId := chi.URLParam(r, "habitId")

	habit, err := hh.HabitService.GetHabitById(ctx, habitId, session.UserId)
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

func (hh *HabitHandler) HandleGetHabitYearDailyCounts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	session := ctx.Value("session").(*models.Session)
	habitId := chi.URLParam(r, "habitId")

	dailyCounts, err := hh.HabitService.GetHabitYearDailyCounts(ctx, habitId, session.UserId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ErrorJSON(w, "Habit Not Found", http.StatusNotFound)
			return
		}
		hh.Logger.Error().Err(err).Msg("HandleGetHabitYearDailyCounts")
		ErrorJSON(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	SendJSON(w, dailyCounts)
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
