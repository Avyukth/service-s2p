package schema

import (
	"context"
	"database/sql"
	_ "embed"
	"errors"
	"fmt"

	"github.com/Avyukth/service3-clone/business/sys/database"
	"github.com/ardanlabs/darwin"
	"github.com/jmoiron/sqlx"
)

var (
	//go:embed sql/schema.sql
	schemaDoc string

	//go:embed sql/seed.sql
	seedDoc string

	//go:embed sql/delete.sql
	deleteDoc string
)

func Migrate(ctx context.Context, db *sqlx.DB) error {

	if err := database.StatusCheck(ctx, db); err != nil {
		return fmt.Errorf("status check database error: %w", err)
	}

	driver, err := darwin.NewGenericDriver(db.DB, darwin.PostgresDialect{})
	if err != nil {
		return err
	}

	d := darwin.New(driver, darwin.ParseMigrations(schemaDoc))
	return d.Migrate()
}

func Seed(ctx context.Context, db *sqlx.DB) error {

	if err := database.StatusCheck(ctx, db); err != nil {
		return fmt.Errorf("status check database error: %w", err)
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	if _, err := tx.Exec(seedDoc); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}
	return tx.Commit()
}

func DeleteAll(ctx context.Context, db *sqlx.DB) error {
	if err := database.StatusCheck(ctx, db); err != nil {
		return fmt.Errorf("status check database error: %w", err)
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() {

		if errTx := tx.Rollback(); errTx != nil {
			if errors.Is(errTx, sql.ErrTxDone) {
				return
			}
			err = fmt.Errorf("rollback error: %w", errTx)
			return
		}
	}()

	if _, err := tx.Exec(deleteDoc); err != nil {
		return fmt.Errorf("exec delete error: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit error: %w", err)
	}
	return nil
}
