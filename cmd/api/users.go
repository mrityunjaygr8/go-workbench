package main

import (
	"encoding/json"
	"net/http"
)

type createUserRequest struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

func (a *application) createUserHandler(w http.ResponseWriter, r *http.Request) {
	var payload createUserRequest
	_ = json.NewDecoder(r.Body).Decode(&payload)
	a.logger.Println(payload)

	userExists, err := a.store.Users.GetUserByEmail(payload.Email)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if userExists != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user := User{
		Email:     payload.Email,
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
		Password:  payload.Password,
	}

	err = a.store.Users.CreateUser(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}
func (a *application) listUsersHandler(w http.ResponseWriter, r *http.Request) {
	users, err := a.store.Users.ListUsers()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}
