package httpkit

import (
	"encoding/json"
	"errors"
	"net/http"
	"route256/cart/internal/domain"
)

var businessErrorStatus = []struct {
	err    error
	status int
}{
	{domain.ErrUserNotFound, http.StatusNotFound},
	{domain.ErrSKUNotFound, http.StatusNotFound},
	{domain.ErrInsufficientStocks, http.StatusConflict},
}

func ClassifyBusinessError(err error) int {
	for _, m := range businessErrorStatus {
		if errors.Is(err, m.err) {
			return m.status
		}
	}
	return http.StatusInternalServerError
}

func WriteError(w http.ResponseWriter, handlerName string, err error, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"handler": handlerName,
		"error":   err.Error(),
	})
}

func WriteJSON(w http.ResponseWriter, handlerName string, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if body == nil {
		return
	}
	if err := json.NewEncoder(w).Encode(body); err != nil {
		WriteError(w, handlerName, err, http.StatusInternalServerError)
	}
}

func DecodeJSON(w http.ResponseWriter, r *http.Request, handlerName string, dst any) bool {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		WriteError(w, handlerName, err, http.StatusBadRequest)
		return false
	}
	return true
}
