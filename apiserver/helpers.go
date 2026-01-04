package apiserver

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
)

type ErrWithStatus struct {
	status int
	err    error
}

func (e *ErrWithStatus) Error() string {
	return e.err.Error()
}

func NewErrWithStatus(status int, err error) *ErrWithStatus {
	return &ErrWithStatus{
		status: status,
		err:    err,
	}
}

func handler(f func(w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			status := http.StatusInternalServerError
			msg := http.StatusText(status)
			if e, ok := err.(*ErrWithStatus); ok {
				status = e.status
				msg = http.StatusText(e.status)
				if status == http.StatusBadRequest || status == http.StatusConflict {
					msg = e.err.Error()
				}
			}
			slog.Error("error executing handler", "error", err, "status", status, "message", msg)
			w.WriteHeader(status)
			if err := json.NewEncoder(w).Encode(ApiResponse[struct{}]{
				Message: msg,
			}); err != nil {
				slog.Error("error encoding response", "error", err)
			}
		}

	})
}

func encode[T any](v T, status int, w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		return fmt.Errorf("failed to encode response: %w", err)
	}
	return nil
}

type Validator interface {
	Validate() error
}

func decode[T Validator](r *http.Request) (T, error) {
	var t T
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		return t, fmt.Errorf("invalid request body: %w", err)
	}

	if err := t.Validate(); err != nil {
		return t, err
	}

	return t, nil
}
