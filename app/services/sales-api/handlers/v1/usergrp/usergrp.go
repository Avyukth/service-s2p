package usergrp

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	userCore "github.com/Avyukth/service3-clone/business/core/user"
	"github.com/Avyukth/service3-clone/business/sys/auth"
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
