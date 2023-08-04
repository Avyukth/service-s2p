package checkgrp

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

type Handlers struct {
	Build string
	Log   *zap.SugaredLogger
}

func (h *Handlers) Readiness(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Status string `json:"status"`
	}{
		Status: "OK",
	}
	statusCode := http.StatusOK
	if err := response(w, statusCode, data); err != nil {
		h.Log.Errorw("readiness", "ERROR", err)
	}

	h.Log.Infow("readiness", "statusCode", statusCode, "method", r.Method, "path", r.URL.Path, "remoteaddr", r.RemoteAddr)
}

func (h *Handlers) Liveness(w http.ResponseWriter, r *http.Request) {

}

func response(w http.ResponseWriter, statusCode int, data interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if _, err := w.Write(jsonData); err != nil {
		return err
	}
	return nil
}
