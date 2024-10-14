package main

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const hmacSampleSecret = "correct-horse-battery-staple"

func (a *application) login(payload loginRequest) (*Token, error) {
	user, err := a.store.GetUserByEmail(payload.Email)
	if err != nil {
		return nil, ErrInvalidCredential
	}

	err = compareHashedPassword(payload.Password, user.Password)
	if err != nil {
		return nil, ErrInvalidCredential
	}

	return a.generateToken(user)
}

func (a *application) revokeToken(token_id uuid.UUID) error {
	token, err := a.store.RetrieveToken(token_id)
	if err != nil {
		return err
	}

	token.Revoked = false
	err = a.store.UpdateToken(token)
	if err != nil {
		return err
	}
	return nil
}

func (a *application) generateToken(user *User) (*Token, error) {
	now := time.Now()
	expiry := now.Add(15 * time.Minute)

	token := Token{
		UserId:    user.ID,
		Expiry:    expiry,
		CreatedAt: now,
		Revoked:   false,
		ID:        uuid.New(),
	}
	tokenTmp := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":      user.Email,
		"first_name": user.FirstName,
		"last_name":  user.LastName,
		"created_at": now,
		"expiry":     expiry,
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := tokenTmp.SignedString([]byte(hmacSampleSecret))
	if err != nil {
		return nil, err
	}
	token.Token = tokenString

	err = a.store.InsertToken(&token)
	if err != nil {
		return nil, err
	}

	return &token, nil
}

func (a *application) listTokens() ([]Token, error) {
	return a.store.ListTokens()
}
