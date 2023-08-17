package user

import (
	"context"
	"fmt"
	"time"

	"github.com/Avyukth/service3-clone/business/sys/auth"
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

func (s Store) Update(ctx context.Context, claims auth.Claims, userID string, uu UpdateUser, now time.Time) error {

	if err := validate.Check(userID); err != nil {
		return database.ErrInvalidID
	}

	if err := validate.Check(uu); err != nil {
		return fmt.Errorf("validating data: %w", err)
	}

	usr, err := s.QueryByID(ctx, claims, userID)
	if err != nil {
		return fmt.Errorf("updating user userID[%s]: %w", userID, err)
	}

	if uu.Name != nil {
		usr.Name = *uu.Name
	}
	if uu.Email != nil {
		usr.Email = *uu.Email
	}
	if uu.Roles != nil {
		usr.Roles = uu.Roles
	}

	if uu.Password != nil {
		pw, err := bcrypt.GenerateFromPassword([]byte(*uu.Password), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("generating password hash: %w", err)
		}
		usr.PasswordHash = pw

	}
	usr.DateUpdated = now
	const q = `
	UPDATE 
		users 
	SET 
		"name"=:name, 
		"email"=:email, 
		"roles"=:roles, 
		"date_updated"=:date_updated 
	WHERE user_id=:user_id`
	if err := database.NamedExecContext(ctx, s.log, s.db, q, usr); err != nil {
		return fmt.Errorf("updating user userID[%s]: %w", userID, err)
	}
	return nil
}

func (s Store) Delete(ctx context.Context, claims auth.Claims, userID string) error {

	if err := validate.Check(userID); err != nil {
		return database.ErrInvalidID
	}

	if !claims.Authorized(auth.RoleAdmin) && claims.Subject != userID {
		return database.ErrForbidden
	}

	data := struct {
		UserID string `db:"user_id"`
	}{
		UserID: userID,
	}

	const q = `
		DELETE FROM 
			users 
		WHERE 
			user_id=:user_id`

	if err := database.NamedExecContext(ctx, s.log, s.db, q, data); err != nil {
		return fmt.Errorf("deleting user userID[%s]: %w", userID, err)
	}

	return nil
}

func (s Store) QueryByID(ctx context.Context, claims auth.Claims, userID string) (User, error) {

	return User{}, nil
}

func (s Store) Query(ctx context.Context, pageNumber int, rowsPerPage int) ([]User, error) {

	data := struct {
		Offset      int `db:"offset"`
		RowsPerPage int `db:"rows_per_page"`
	}{
		Offset:      (pageNumber - 1) * rowsPerPage,
		RowsPerPage: rowsPerPage,
	}

	const q = `
	SELECT
		*
	FROM
		users
	ORDER BY
		user_id
	OFFSET :offset ROWS FETCH NEXT :rows_per_page ROWS ONLY`

	var users []User
	if err := database.NamedQuerySlice(ctx, s.log, s.db, q, data, &users); err != nil {
		if err == database.ErrNotFound {
			return nil, database.ErrNotFound
		}
		return nil, fmt.Errorf("selecting users: %w", err)
	}

	return users, nil

}
