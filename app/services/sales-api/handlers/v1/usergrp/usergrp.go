package usergrp

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	userCore "github.com/Avyukth/service3-clone/business/core/user"
	"github.com/Avyukth/service3-clone/business/data/store/user"
	"github.com/Avyukth/service3-clone/business/sys/auth"
	"github.com/Avyukth/service3-clone/business/sys/database"
	"github.com/Avyukth/service3-clone/business/sys/validate"

	"github.com/Avyukth/service3-clone/foundation/web"
)

type Handlers struct {
	User userCore.Core
	Auth *auth.Auth
}

func (h Handlers) Query(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	page := web.Param(r, "page")
	pageNumber, err := strconv.Atoi(page)

	if err != nil {
		return validate.NewRequestError(fmt.Errorf("invalid page format: [%s]", page), http.StatusBadRequest)
	}
	rows := web.Param(r, "rows")
	rowsPerPage, err := strconv.Atoi(rows)

	if err != nil {
		return validate.NewRequestError(fmt.Errorf("invalid rows format [%s]", rows), http.StatusBadRequest)
	}
	users, err := h.User.Query(ctx, pageNumber, rowsPerPage)
	if err != nil {
		return fmt.Errorf("unable to query for users : %w", err)
	}

	return web.Respond(ctx, w, users, http.StatusOK)
}

func (h Handlers) QueryByID(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	claims, err := auth.GetClaims(ctx)

	if err != nil {
		return errors.New("claims missing in context")
	}

	id := web.Param(r, "id")

	usr, err := h.User.QueryById(ctx, claims, id)
	if err != nil {
		switch validate.Cause(err) {
		case database.ErrInvalidID:
			return validate.NewRequestError(err, http.StatusBadRequest)
		case database.ErrNotFound:
			return validate.NewRequestError(err, http.StatusNotFound)
		case database.ErrForbidden:
			return validate.NewRequestError(err, http.StatusForbidden)
		default:
			return fmt.Errorf("ID[%s]: %w", id, err)
		}

	}

	return web.Respond(ctx, w, usr, http.StatusOK)
}

func (h Handlers) Create(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	v, err := web.GetValues(ctx)
	if err != nil {
		return web.NewShutdownError("web value missing from context")
	}

	var nu user.NewUser

	if err := web.Decode(r, &nu); err != nil {
		return fmt.Errorf("unable to decode payload: %w", err)
	}
	usr, err := h.User.Create(ctx, nu, v.Now)
	if err != nil {
		return fmt.Errorf("unable to create user[%+v]: %w", &usr, err)
	}

	return web.Respond(ctx, w, usr, http.StatusCreated)
}

func (h Handlers) Update(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	v, err := web.GetValues(ctx)
	if err != nil {
		return web.NewShutdownError("web value missing from context")
	}
	claims, err := auth.GetClaims(ctx)
	if err != nil {
		return errors.New("claims missing in context")
	}

	var upd user.UpdateUser
	if err := web.Decode(r, &upd); err != nil {
		return fmt.Errorf("unable to decode payload: %w", err)
	}
	id := web.Param(r, "id")
	if err := h.User.Update(ctx, claims, id, upd, v.Now); err != nil {
		switch validate.Cause(err) {
		case database.ErrInvalidID:
			return validate.NewRequestError(err, http.StatusBadRequest)
		case database.ErrNotFound:
			return validate.NewRequestError(err, http.StatusNotFound)
		case database.ErrForbidden:
			return validate.NewRequestError(err, http.StatusForbidden)
		default:
			return fmt.Errorf("ID[%s] User[%+v]: %w", id, &upd, err)
		}
	}

	return web.Respond(ctx, w, nil, http.StatusNoContent)

}

func (h Handlers) Delete(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	claims, err := auth.GetClaims(ctx)
	if err != nil {
		return errors.New("claims missing in context")
	}

	id := web.Param(r, "id")

	if err := h.User.Delete(ctx, claims, id); err != nil {
		switch validate.Cause(err) {
		case database.ErrInvalidID:
			return validate.NewRequestError(err, http.StatusBadRequest)
		case database.ErrNotFound:
			return validate.NewRequestError(err, http.StatusNotFound)
		case database.ErrForbidden:
			return validate.NewRequestError(err, http.StatusForbidden)
		default:
			return fmt.Errorf("ID[%s]: %w", id, err)
		}
	}
	return web.Respond(ctx, w, nil, http.StatusNoContent)
}

func (h Handlers) Token(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	v, err := web.GetValues(ctx)
	if err != nil {
		return web.NewShutdownError("web value missing from context")
	}

	email, pass, ok := r.BasicAuth()
	if !ok {
		err := errors.New("must provide email and password in Basic auth")
		return validate.NewRequestError(err, http.StatusUnauthorized)
	}

	claims, err := h.User.Authenticate(ctx, v.Now, email, pass)
	if err != nil {
		switch validate.Cause(err) {
		case database.ErrNotFound:
			return validate.NewRequestError(err, http.StatusNotFound)
		case database.ErrAuthenticationFailure:
			return validate.NewRequestError(err, http.StatusUnauthorized)
		default:
			return fmt.Errorf("email[%s]: %w", email, err)
		}
	}

	var tkn struct {
		Token string `json:"token"`
	}

	tkn.Token, err = h.Auth.GenerateToken(claims)

	if err != nil {
		return fmt.Errorf("unable to generate token: %w", err)
	}
	return web.Respond(ctx, w, tkn, http.StatusOK)
}
