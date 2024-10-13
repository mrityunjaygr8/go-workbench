package main

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type store interface {
	UserManager
	StoreMeta
	AuthManager
}

type StoreMeta interface {
	Persist() error
	Load() error
}

type User struct {
	Email     string    `json:"email"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Password  string    `json:"-"`
	ID        uuid.UUID `json:"id"`
}

type Token struct {
	Token     string    `json:"token"`
	UserId    uuid.UUID `json:"user_id"`
	Revoked   bool      `json:"revoked"`
	Expiry    time.Time `json:"expiry"`
	CreatedAt time.Time `json:"created_at"`
}

type UserManager interface {
	ListUsers() ([]User, error)
	InsertUser(user User) error
	GetUserByEmail(email string) (*User, error)
}

type AuthManager interface {
	InsertToken(token Token) error
	ListTokens() ([]Token, error)
	// RevokeToken(token string) error
	// Logout() error
}

var ErrInvalidCredential = fmt.Errorf("invalid email or password provided")
var ErrUserNotFound = fmt.Errorf("user not found")
