package usergrp

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	userCore "github.com/Avyukth/service3-clone/business/core/user"
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
