package services

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/Fozzyack/habit-tracker/internal/models"
	"github.com/Fozzyack/habit-tracker/internal/store"
)

var ErrInvalidRecordHabitDate = errors.New("invalid record habit date")

type HabitService struct {
	TxManager            TxManager
	HabitStore           store.HabitStore
	HabitDailyTotalStore store.HabitDailyTotalsStore
	HabitLogStore        store.HabitLogStore
}

func NewHabitService(txManager TxManager, habitStore store.HabitStore, habitDailyTotalStore store.HabitDailyTotalsStore, habitStoreLog store.HabitLogStore) *HabitService {
	return &HabitService{
		TxManager:            txManager,
		HabitStore:           habitStore,
		HabitDailyTotalStore: habitDailyTotalStore,
		HabitLogStore:        habitStoreLog,
	}
}

func (hs *HabitService) RecordHabit(ctx context.Context, userId string, habitTotalReq *models.RecordHabitRequest) (*models.HabitDailyTotal, error) {
	date, err := time.Parse("2006-01-02", habitTotalReq.Date)
	if err != nil {
		return nil, ErrInvalidRecordHabitDate
	}

	var habitTotal *models.HabitDailyTotal
	err = hs.TxManager.WithTx(ctx, func(tx *sql.Tx) error {
		var err error

		habitTotal, err = hs.HabitDailyTotalStore.GetDailyHabitTotalTx(ctx, tx, habitTotalReq.HabitId, userId, date)
		if errors.Is(err, sql.ErrNoRows) {
			habitTotal, err = hs.HabitDailyTotalStore.CreateDailyHabitTotal(ctx, tx, habitTotalReq.Amount, userId, habitTotalReq.HabitId, date)
		} else if err != nil {
			return err
		} else {
			habitTotal.Amount += habitTotalReq.Amount
			habitTotal, err = hs.HabitDailyTotalStore.UpdateDailyHabitTotal(ctx, tx, habitTotal)
		}
		if err != nil {
			return err
		}

		_, err = hs.HabitLogStore.CreateHabitLog(ctx, tx, habitTotalReq.Amount, userId, habitTotalReq.HabitId)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}
	return habitTotal, nil

}

func (hs *HabitService) CreateHabit(ctx context.Context, userId string, habitReq *models.NewHabitRequest) (*models.Habit, error) {
	if habitReq.Increment == 0 {
		habitReq.Increment = 1
	}

	var habit *models.Habit
	err := hs.TxManager.WithTx(ctx, func(tx *sql.Tx) error {
		var err error
		habit, err = hs.HabitStore.CreateHabit(ctx, tx, userId, habitReq)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return habit, nil
}

func (hs *HabitService) GetHabitById(ctx context.Context, id, userId string) (*models.Habit, error) {
	habit, err := hs.HabitStore.GetHabitById(ctx, id, userId)
	if err != nil {
		return nil, err
	}

	return habit, nil
}

func (hs *HabitService) UpdateHabit(ctx context.Context, id, userId string, habitReq *models.UpdateHabitRequest) (*models.Habit, error) {
	var habit *models.Habit
	err := hs.TxManager.WithTx(ctx, func(tx *sql.Tx) error {
		var err error
		habit, err = hs.HabitStore.UpdateHabit(ctx, tx, id, userId, habitReq)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return habit, nil
}

func (hs *HabitService) GetHabitsByUserId(ctx context.Context, userId string) ([]*models.Habit, error) {
	habits, err := hs.HabitStore.GetHabitsByUserId(ctx, userId)
	if err != nil {
		return nil, err
	}

	return habits, nil
}

func (hs *HabitService) GetHabitYearDailyCounts(ctx context.Context, habitId, userId string) ([]*models.HabitDailyCount, error) {
	_, err := hs.HabitStore.GetHabitById(ctx, habitId, userId)
	if err != nil {
		return nil, err
	}

	dailyCounts, err := hs.HabitDailyTotalStore.GetHabitYearDailyCounts(ctx, habitId, userId)
	if err != nil {
		return nil, err
	}

	return dailyCounts, nil
}
