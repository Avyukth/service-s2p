package web

import (
	"context"
	"net/http"
)

func Respond(ctx context.Context, w http.ResponseWriter, data interface{}, statusCode int) error {

	if statusCode == http.StatusNoContent {
		w.WriteHeader(statusCode)
		return nil
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	w.Header().Set("Content-type", "application/json")

	w.WriteHeader(statusCode)
	if _, err := w.Write(jsonData); err != nil {
		return err
	}

	return nil
}
