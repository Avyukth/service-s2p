package mid

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Avyukth/service3-clone/foundation/web"
)

func Panics() web.Middleware {
	m := func(handler web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) (err error) {
			defer func() {

				if rec := recover(); rec != nil {
					err = fmt.Errorf("PANIC [%v]", rec)
				}
			}()
			return handler(ctx, w, r)
		}
		return h
	}

	return m
}
