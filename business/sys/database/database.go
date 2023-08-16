package database

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"strings"
	"time"

	"github.com/Avyukth/service3-clone/foundation/web"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

var (
	ErrorNotFound            = errors.New("not found")
	ErrInvalidID             = errors.New("ID is not in proper form")
	ErrAuthenticationFailure = errors.New("authentication failed")
	ErrForbidden             = errors.New("attempt action not allowed")
)

type Config struct {
	User         string
	Password     string
	Host         string
	Name         string
	MaxIdleConns int
	MaxOpenConns int
	DisableTLS   bool
}

func Open(cfg Config) (*sqlx.DB, error) {
	sslMode := "require"

	if cfg.DisableTLS {
		sslMode = "disable"
	}

	q := make(url.Values)
	q.Set("sslmode", sslMode)
	q.Set("timezone", "UTC")

	u := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(cfg.User, cfg.Password),
		Host:     cfg.Host,
		Path:     cfg.Name,
		RawQuery: q.Encode(),
	}
	fmt.Println("URL string ......................: ", url.UserPassword(cfg.User, cfg.Password), u.String())
	db, err := sqlx.Open("postgres", u.String())
	if err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	fmt.Println("DB connection Successfully Open", db.Stats())
	return db, nil
}

func StatusCheck(ctx context.Context, db *sqlx.DB) error {
	var (
		deadline time.Time
		ok       bool
	)
	if deadline, ok = ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, time.Second)
		defer cancel()
	}
	currrentTime := time.Now()
	diff := currrentTime.Sub(deadline)
	fmt.Println("DB connection StatusCheck deadline .............", deadline, diff)

	var pingError error
	for attempts := 1; ; attempts++ {
		pingError = db.Ping()

		if pingError == nil {
			break
		}

		time.Sleep(time.Duration(attempts) * 10 * time.Millisecond)

		// if ctx.Err() != nil {
		// 	return ctx.Err()
		// }
		if err := ctx.Err(); err == context.DeadlineExceeded {
			return fmt.Errorf("database is unreachable, deadline exceeded after retries, pingError: %w", pingError)
		} else if err == context.Canceled {
			return fmt.Errorf("database check was cancelled")
		}

	}

	if ctx.Err() != nil {
		return ctx.Err()
	}

	// Run a simple query to determine connectivity.
	// Running this query forces a round trip through the database.
	const q = `SELECT true`
	var tmp bool
	return db.QueryRowContext(ctx, q).Scan(&tmp)
}

func NameExecContext(ctx context.Context, log *zap.SugaredLogger, db *sqlx.DB, query string, data interface{}, dest interface{}) error {

	q := queryString(query, data)
	log.Infow("database.NameExecContext", "traceid", web.GetTraceID(ctx), "query", q)

	if _, err := db.ExecContext(ctx, q); err != nil {
		return err
	}
	return nil
}

func NamedQuerySlice(ctx context.Context, log *zap.SugaredLogger, db *sqlx.DB, query string, data interface{}, dest any) error {

	q := queryString(query, data)
	log.Infow("database.NamedQuerySlice", "traceid", web.GetTraceID(ctx), "query", q)

	val := reflect.ValueOf(dest)

	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Slice {
		return errors.New("must provide a pointer to a slice")
	}

	rows, err := db.NamedQueryContext(ctx, query, data)
	if err != nil {
		return err
	}
	slice := val.Elem()
	for rows.Next() {
		v := reflect.New(slice.Type().Elem())
		if err := rows.StructScan(v.Interface()); err != nil {
			return err
		}
		slice.Set(reflect.Append(slice, v.Elem()))
	}
	return nil
}

func NamedQueryStruct(ctx context.Context, log *zap.SugaredLogger, db *sqlx.DB, query string, data interface{}, dest any) error {

	q := queryString(query, data)
	log.Infow("database.NamedQueryStruct", "traceid", web.GetTraceID(ctx), "query", q)

	rows, err := db.NamedQueryContext(ctx, query, data)

	if err != nil {
		return err
	}

	if !rows.Next() {
		return ErrorNotFound
	}

	if err := rows.StructScan(dest); err != nil {
		return err
	}

	return nil
}

func queryString(query string, args ...interface{}) string {

	query, params, err := sqlx.Named(query, args)
	if err != nil {
		return err.Error()
	}

	for _, param := range params {
		var value string
		switch v := param.(type) {
		case string:
			value = fmt.Sprintf("%q", v)
		case []byte:
			value = fmt.Sprintf("%q", string(v))

		default:
			value = fmt.Sprintf("%v", v)
		}
		query = strings.Replace(query, "?", value, 1)
	}

	query = strings.ReplaceAll(query, "\t", "")
	query = strings.ReplaceAll(query, "\n", "")

	return strings.Trim(query, " ")
}
