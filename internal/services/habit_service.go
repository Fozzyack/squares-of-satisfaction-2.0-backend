package services

import (
	"context"
	"database/sql"

	"github.com/Fozzyack/habit-tracker/internal/models"
	"github.com/Fozzyack/habit-tracker/internal/store"
)

type HabitService struct {
	TxManager  TxManager
	HabitStore store.HabitStore
}

func NewHabitService(txManager TxManager, habitStore store.HabitStore) *HabitService {
	return &HabitService{
		TxManager:  txManager,
		HabitStore: habitStore,
	}
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

func (hs *HabitService) GetHabitById(id string) (*models.Habit, error) {
	habit, err := hs.HabitStore.GetHabitById(id)
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
