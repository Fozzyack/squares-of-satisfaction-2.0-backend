package services

import (
	"context"
	"database/sql"
)



type TxManager interface {
	WithTx(ctx context.Context, fn func(tx *sql.Tx) error) error
}

type SQLTxManager struct {
	db *sql.DB
}

func NewSQLTxManager(db *sql.DB) TxManager {
	return &SQLTxManager{db: db}
}

func (stm *SQLTxManager) WithTx(ctx context.Context, fn func(tx *sql.Tx) error) error {

	tx, err := stm.db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	err = fn(tx)
	if err != nil {
		return err 
	}

	err = tx.Commit() 
	if err != nil {
		return err
	}

	return nil
}

