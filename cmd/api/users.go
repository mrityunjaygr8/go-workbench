package main

import (
	"encoding/json"
	"errors"
	"net/http"
)

type createUserRequest struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type ErrorMessage struct {
	Error string `json:"error"`
}

func (a *application) handleUserCreate(w http.ResponseWriter, r *http.Request) {
	var payload createUserRequest
	_ = json.NewDecoder(r.Body).Decode(&payload)

	err := a.createUser(payload)
	if err != nil {
		if errors.Is(err, ErrUserExists) {
			msg := ErrorMessage{
				Error: err.Error(),
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(msg)
			return

		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}
func (a *application) handleUserList(w http.ResponseWriter, r *http.Request) {
	users, err := a.listUsers()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}
