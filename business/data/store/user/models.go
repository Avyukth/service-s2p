package user

import (
	"time"

	"github.com/Avyukth/service3-clone/business/sys/auth"
)

type User struct {
	ID           string    `db:"user_id" json:"id"`
	Name         string    `db:"name" json:"name"`
	Email        string    `db:"email" json:"email"`
	Roles        string    `db:"roles" json:"roles"`
	PasswordHash []byte    `db:"password_hash" json:"-"`
	DateCreated  time.Time `db:"date_created" json:"date_created"`
	DateUpdated  time.Time `db:"date_updated" json:"date_updated"`
}

type NewUser struct {
	Name            string      `json:"name" validate:"required"`
	Email           string      `json:"email" validate:"required,email"`
	Roles           []auth.Role `json:"roles" validate:"required"`
	Password        string      `json:"password" validate:"required"`
	PasswordConfirm string      `json:"password_confirm" validate:"eqfield=Password"`
}

type UpdateUser struct {
	Name            *string     `json:"name"`
	Email           *string     `json:"email" validate:"omitempty,email"`
	Roles           []auth.Role `json:"roles"`
	Password        *string     `json:"password"`
	PasswordConfirm *string     `json:"password_confirm" validate:"omitempty,eqfield=Password"`
}
