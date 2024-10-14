package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (a *application) handleLogin(w http.ResponseWriter, r *http.Request) {
	var payload loginRequest
	_ = json.NewDecoder(r.Body).Decode(&payload)

	token, err := a.login(payload)
	if err != nil {
		if errors.Is(err, ErrInvalidCredential) {
			msg := ErrorMessage{
				Error: err.Error(),
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(msg)
			return
		}
		a.logger.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(token)
	return
}

func (a *application) handleTokenList(w http.ResponseWriter, r *http.Request) {
	tokens, err := a.listTokens()
	if err != nil {
		a.logger.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(tokens)
	return

}

type tokenRevokeRequest struct {
	TokenID string `json:"token_id"`
}

func (a *application) handleTokenRevoke(w http.ResponseWriter, r *http.Request) {
	var payload tokenRevokeRequest
	_ = json.NewDecoder(r.Body).Decode(&payload)

	token_id, err := uuid.Parse(payload.TokenID)
	if err != nil {
		msg := ErrorMessage{
			Error: fmt.Errorf("not a valid uuid").Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(msg)
		return
	}
	err = a.revokeToken(token_id)

	if err != nil {
		if errors.Is(err, ErrNotFound) {
			msg := ErrorMessage{
				Error: fmt.Errorf("not found").Error(),
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(msg)
			return

		}
		a.logger.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return

	}

	w.WriteHeader(http.StatusOK)
	return

}
