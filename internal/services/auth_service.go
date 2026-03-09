package services

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/Fozzyack/habit-tracker/internal/auth"
	"github.com/Fozzyack/habit-tracker/internal/models"
	"github.com/Fozzyack/habit-tracker/internal/store"
	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidCredentials = errors.New("invalid credentials")

type AuthService struct {
	TxManager    TxManager
	UserStore    store.UserStore
	SessionStore store.SessionStore
}

func NewAuthService(txManager TxManager, userStore store.UserStore, sessionStore store.SessionStore) *AuthService {
	return &AuthService{
		TxManager:    txManager,
		UserStore:    userStore,
		SessionStore: sessionStore,
	}
}

func (as *AuthService) CreateNewUser(ctx context.Context, userReq *models.NewUserRequest) (*models.User, *models.Session, error) {

	var user *models.User
	var session *models.Session

	token, err := auth.GenerateToken()
	if err != nil {
		return nil, nil, err
	}
	password, err := auth.HashPassword(userReq.Password)
	if err != nil {
		return nil, nil, err
	}

	err = as.TxManager.WithTx(ctx, func(tx *sql.Tx) error {
		var err error
		user, err = as.UserStore.CreateUser(ctx, tx, password, userReq)
		if err != nil {
			return err
		}
		session, err = as.SessionStore.CreateSession(ctx, tx, user.Id, token, time.Now().UTC().Add(time.Hour*24))
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, nil, err
	}
	return user, session, nil
}

func (as *AuthService) LoginUser(ctx context.Context, loginReq *models.LoginUserRequest) (*models.Session, error) {
	user, err := as.UserStore.GetUserByEmail(ctx, loginReq.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	err = auth.VerifyPassword(user.PasswordHash, loginReq.Password)
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	token, err := auth.GenerateToken()
	if err != nil {
		return nil, err
	}

	var session *models.Session
	err = as.TxManager.WithTx(ctx, func(tx *sql.Tx) error {
		session, err = as.SessionStore.CreateSession(ctx, tx, user.Id, token, time.Now().UTC().Add(time.Hour*24))
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return session, nil
}
