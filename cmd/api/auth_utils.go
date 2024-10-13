package main

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

func hashPassword(raw string) (string, error) {

	password, err := bcrypt.GenerateFromPassword([]byte(raw), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(password), nil
}

func compareHashedPassword(raw, actual string) error {
	err := bcrypt.CompareHashAndPassword([]byte(actual), []byte(raw))
	if err != nil {
		return err
	}

	return nil
}

func createToken(user User) (Token, error) {
	now := time.Now()
	token := Token{UserId: user.ID, CreatedAt: now, Expiry: now.Add(5 * time.Minute)}
	// token_string :=

	return token, nil
}
