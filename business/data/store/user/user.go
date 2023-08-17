package user

import (
	"context"
	"fmt"
	"time"

	"github.com/Avyukth/service3-clone/business/sys/database"
	"github.com/Avyukth/service3-clone/business/sys/validate"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
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

func (s Store) Create(ctx context.Context, nu *NewUser, now time.Time) (User, error) {

	if err := validate.Check(nu); err != nil {
		return User{}, fmt.Errorf("validating data: %w", err)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(nu.Password), bcrypt.DefaultCost)

	if err != nil {
		return User{}, fmt.Errorf("generating password hash: %w", err)
	}

	usr := User{
		ID:           validate.GenerateID(),
		Name:         nu.Name,
		Email:        nu.Email,
		PasswordHash: hash,
		DateCreated:  now,
		DateUpdated:  now,
	}

	const q = `INSERT INTO users (id, name, email, password_hash, date_created, date_updated) VALUES (:user_id, :name, :email, :password_hash, :date_created, :date_updated)`

	if err := database.NamedExecContext(ctx, s.log, s.db, q, usr); err != nil {
		return User{}, fmt.Errorf("inserting user: %w", err)
	}
	return usr, nil
}
