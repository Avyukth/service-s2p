package user

import (
	"context"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type Store struct {
	log *zap.SugaredLogger
	db  *sqlx.DB
}

func NewStore(log *zap.SugaredLogger, db *sqlx.DB) *Store {
	return &Store{
		log: log,
		db:  db,
	}
}

func (s Store) Create(ctx context.Context, nu *NewUser) error {

	return nil
}
