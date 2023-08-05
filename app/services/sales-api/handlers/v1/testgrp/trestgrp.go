package testgrp

import (
	"context"
	"net/http"

	"github.com/Avyukth/service3-clone/foundation/web"
	"go.uber.org/zap"
)

type Handlers struct {
	Log *zap.SugaredLogger
}

func (h Handlers) Test(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	status := struct {
		Status string `json:"status"`
	}{
		Status: "OK",
	}
	statusCode := http.StatusOK

	return web.Respond(ctx, w, status, statusCode)

}
