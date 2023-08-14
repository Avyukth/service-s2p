package mid

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/Avyukth/service3-clone/business/sys/auth"
	"github.com/Avyukth/service3-clone/business/sys/validate"
	"github.com/Avyukth/service3-clone/foundation/web"
)

func Authenticate(a *auth.Auth) web.Middleware {
	m := func(handler web.Handler) web.Handler {

		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			authString := r.Header.Get("authorization")

			parts := strings.Split(authString, " ")

			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				err := errors.New("expected authorization header to be Bearer token format: bearer <token>")

				return validate.NewRequestError(err, http.StatusUnauthorized)
			}

			claims, err := a.ValidateToken(parts[1])
			if err != nil {
				return validate.NewRequestError(err, http.StatusUnauthorized)
			}

			ctx = auth.SetClaims(ctx, claims)

			return handler(ctx, w, r)
		}
		return h
	}

	return m
}

func Authorize(roles ...string) web.Middleware {
	m := func(handler web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

			claims, err := auth.GetClaims(ctx)
			if err != nil {
				return validate.NewRequestError(fmt.Errorf("you are not authorized for that action, no claims"), http.StatusForbidden)
			}

			if !claims.Authorized(roles...) {

				return validate.NewRequestError(fmt.Errorf("you are not authorized for that action, claims[%v] roles[%v]", claims.Roles, roles), http.StatusForbidden)
			}

			return handler(ctx, w, r)
		}

		return h
	}
	return m
}
