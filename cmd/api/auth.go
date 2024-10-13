package main

import (
	"encoding/json"
	"errors"
	"net/http"
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

func (a *application) handleListTokens(w http.ResponseWriter, r *http.Request) {
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
