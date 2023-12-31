package user

import (
	"context"
	"fmt"
	"time"

	"github.com/Avyukth/service3-clone/business/data/store/user"
	"github.com/Avyukth/service3-clone/business/sys/auth"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type Core struct {
	log  *zap.SugaredLogger
	user user.Store
}

func NewCore(log *zap.SugaredLogger, db *sqlx.DB) Core {
	return Core{
		log:  log,
		user: user.NewStore(log, db),
	}
}

func (c Core) Create(ctx context.Context, nu user.NewUser, now time.Time) (user.User, error) {

	// PERFORM PRE BUSINESSES OPERATIONS

	usr, err := c.user.Create(ctx, nu, now)
	if err != nil {
		return user.User{}, fmt.Errorf("create user failed: %w", err)
	}

	// PERFORM POST BUSINESSES OPERATIONS

	return usr, nil
}

func (c Core) Update(ctx context.Context, claims auth.Claims, userID string, uu user.UpdateUser, now time.Time) error {
	// PERFORM PRE BUSINESSES OPERATIONS

	if err := c.user.Update(ctx, claims, userID, uu, now); err != nil {
		return fmt.Errorf("update user failed: %w", err)
	}

	// PERFORM POST BUSINESSES OPERATIONS

	return nil

}

func (c Core) Delete(ctx context.Context, claims auth.Claims, userID string) error {
	// PERFORM PRE BUSINESSES OPERATIONS

	if err := c.user.Delete(ctx, claims, userID); err != nil {
		return fmt.Errorf("delete user failed: %w", err)
	}

	// PERFORM POST BUSINESSES OPERATIONS

	return nil

}

func (c Core) Query(ctx context.Context, pageNumber int, rowsPerPage int) ([]user.User, error) {
	// PERFORM PRE BUSINESSES OPERATIONS

	users, err := c.user.Query(ctx, pageNumber, rowsPerPage)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}

	// PERFORM POST BUSINESSES OPERATIONS

	return users, nil

}

func (c Core) QueryById(ctx context.Context, claims auth.Claims, userId string) (user.User, error) {
	// PERFORM PRE BUSINESSES OPERATIONS

	usr, err := c.user.QueryByID(ctx, claims, userId)
	if err != nil {
		return user.User{}, fmt.Errorf("query user by id failed: %w", err)
	}

	// PERFORM POST BUSINESSES OPERATIONS

	return usr, nil

}

func (c Core) QueryByEmail(ctx context.Context, claims auth.Claims, email string) (user.User, error) {
	// PERFORM PRE BUSINESSES OPERATIONS

	usr, err := c.user.QueryByEmail(ctx, claims, email)
	if err != nil {
		return user.User{}, fmt.Errorf("update user by email failed: %w", err)
	}

	// PERFORM POST BUSINESSES OPERATIONS

	return usr, nil
}

func (c Core) Authenticate(ctx context.Context, now time.Time, email string, password string) (auth.Claims, error) {
	// PERFORM PRE BUSINESSES OPERATIONS

	claims, err := c.user.Authenticate(ctx, now, email, password)
	if err != nil {
		return auth.Claims{}, fmt.Errorf("authentication fail: %w", err)
	}

	// PERFORM POST BUSINESSES OPERATIONS

	return claims, nil
}
