package database

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
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

func Open(c Config) (*sqlx.DB, error) {

}
