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

type CreatedAt struct {
	CreatedAt time.Time `json:"created_at"`
}

type User struct {
	Email     string    `json:"email"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Password  string    `json:"-"`
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
}

// func (u *User) MarshalJSON() ([]byte, error) {
// 	type Alias User
// 	return json.Marshal(&struct {
// 		CreatedAt int64 `json:"created_at"`
// 		*Alias
// 	}{
// 		CreatedAt: u.CreatedAt.Unix(),
// 		Alias:     (*Alias)(u),
// 	})
// }
// func (u *User) UnmarshalJSON(data []byte) error {
// 	type Alias User
// 	aux := &struct {
// 		CreatedAt int64 `json:"lastSeen"`
// 		*Alias
// 	}{
// 		Alias: (*Alias)(u),
// 	}
// 	if err := json.Unmarshal(data, &aux); err != nil {
// 		return err
// 	}
// 	u.CreatedAt = time.Unix(aux.CreatedAt, 0)
// 	return nil
// }

type Token struct {
	Token     string    `json:"token"`
	UserId    uuid.UUID `json:"user_id"`
	Revoked   bool      `json:"revoked"`
	Expiry    time.Time `json:"expiry"`
	CreatedAt time.Time `json:"created_at"`
	ID        uuid.UUID `json:"id"`
}

// func (u *Token) MarshalJSON() ([]byte, error) {
// 	type Alias Token
// 	return json.Marshal(&struct {
// 		CreatedAt int64 `json:"created_at"`
// 		Expiry    int64 `json:"expiry"`
// 		*Alias
// 	}{
// 		CreatedAt: u.CreatedAt.Unix(),
// 		Expiry:    u.Expiry.Unix(),
// 		Alias:     (*Alias)(u),
// 	})
// }
// func (u *Token) UnmarshalJSON(data []byte) error {
// 	type Alias Token
// 	aux := &struct {
// 		CreatedAt int64 `json:"created_at"`
// 		Expiry    int64 `json:"expiry"`
// 		*Alias
// 	}{
// 		Alias: (*Alias)(u),
// 	}
// 	if err := json.Unmarshal(data, &aux); err != nil {
// 		return err
// 	}
// 	u.CreatedAt = time.Unix(aux.CreatedAt, 0)
// 	u.Expiry = time.Unix(aux.Expiry, 0)
// 	return nil
// }

type UserManager interface {
	ListUsers() ([]User, error)
	InsertUser(*User) error
	GetUserByEmail(string) (*User, error)
}

type AuthManager interface {
	InsertToken(*Token) error
	ListTokens() ([]Token, error)
	UpdateToken(*Token) error
	RetrieveToken(uuid.UUID) (*Token, error)
	// RevokeToken(token string) error
	// Logout() error
}

var ErrInvalidCredential = fmt.Errorf("invalid email or password provided")
var ErrNotFound = fmt.Errorf("object not found")
