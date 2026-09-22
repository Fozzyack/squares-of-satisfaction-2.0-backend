package services

import (
	"context"
	"database/sql"

	"github.com/Fozzyack/habit-tracker/internal/store"
)

type UserService struct {
	TxManager            TxManager
	UserStore            store.UserStore
	SessionStore         store.SessionStore
	HabitStore           store.HabitStore
	HabitDailyTotalStore store.HabitDailyTotalsStore
	HabitLogStore        store.HabitLogStore
}

func NewUserService(
	txManager TxManager,
	userStore store.UserStore,
	sessionStore store.SessionStore,
	habitStore store.HabitStore,
	habitDailyTotalStore store.HabitDailyTotalsStore,
	habitLogStore store.HabitLogStore,
) *UserService {
	return &UserService{
		TxManager:            txManager,
		UserStore:            userStore,
		SessionStore:         sessionStore,
		HabitStore:           habitStore,
		HabitDailyTotalStore: habitDailyTotalStore,
		HabitLogStore:        habitLogStore,
	}
}

func (us *UserService) DeleteUser(ctx context.Context, userId string) error {
	return us.TxManager.WithTx(ctx, func(tx *sql.Tx) error {
		if err := us.HabitLogStore.DeleteHabitLogsByUserIdTx(ctx, tx, userId); err != nil {
			return err
		}

		if err := us.HabitDailyTotalStore.DeleteDailyHabitTotalsByUserIdTx(ctx, tx, userId); err != nil {
			return err
		}

		if err := us.HabitStore.DeleteHabitsByUserIdTx(ctx, tx, userId); err != nil {
			return err
		}

		if err := us.SessionStore.DeleteSessionsByUserIdTx(ctx, tx, userId); err != nil {
			return err
		}

		return us.UserStore.DeleteUserTx(ctx, tx, userId)
	})
}
