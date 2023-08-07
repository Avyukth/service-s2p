package testgrp

import (
	"context"
	"math/rand"
	"net/http"

	"github.com/Avyukth/service3-clone/business/sys/validate"
	"github.com/Avyukth/service3-clone/foundation/web"
	"go.uber.org/zap"
)

type Handlers struct {
	Log *zap.SugaredLogger
}

func (h Handlers) Test(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	if n := rand.Intn(100); n%2 == 0 {
		return validate.NewRequestError(errors.New("untrusted error"), http.StatusBadRequest)
	}

	status := struct {
		Status string `json:"status"`
	}{
		Status: "OK",
	}
	statusCode := http.StatusOK

	return web.Respond(ctx, w, status, statusCode)

}
