package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

var ErrUserExists = fmt.Errorf("user already exists")

func (a *application) createUser(payload createUserRequest) error {
	_, err := a.store.GetUserByEmail(payload.Email)
	if !errors.Is(err, ErrNotFound) || err == nil {
		return ErrUserExists
	}

	user := User{
		Email:     payload.Email,
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
		Password:  payload.Password,
		CreatedAt: time.Now(),
	}
	user.ID = uuid.New()
	user.Password, err = hashPassword(user.Password)
	if err != nil {
		return err
	}
	err = a.store.InsertUser(&user)
	return err
}

func (a *application) getUserByEmail(email string) (*User, error) {
	return a.store.GetUserByEmail(email)
}

func (a *application) listUsers() ([]User, error) {
	return a.store.ListUsers()
}
