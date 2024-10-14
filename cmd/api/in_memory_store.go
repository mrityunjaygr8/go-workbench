package main

import (
	"encoding/json"
	"os"

	"github.com/google/uuid"
)

const FILE_NAME = "store.json"

func NewInMemoryStore() (store, error) {
	s := &InMemoryStore{
		users:  make([]User, 0),
		tokens: make([]Token, 0),
	}

	// Try to load existing data
	err := s.Load()
	if err != nil {
		// If the file doesn't exist, it's not an error, we'll start with an empty store
		if !os.IsNotExist(err) {
			return nil, err
		}
	}

	return s, nil
}

type InMemoryStore struct {
	users  []User
	tokens []Token
}

func (u *InMemoryStore) ListUsers() ([]User, error) {
	return u.users, nil
}

func (u *InMemoryStore) InsertUser(user *User) error {
	u.users = append(u.users, *user)
	return nil
}

func (u *InMemoryStore) GetUserByEmail(email string) (*User, error) {
	for _, user := range u.users {
		if user.Email == email {
			return &user, nil
		}
	}
	return nil, ErrNotFound
}

func (u *InMemoryStore) InsertToken(token *Token) error {
	u.tokens = append(u.tokens, *token)
	return nil

}

func (u *InMemoryStore) ListTokens() ([]Token, error) {
	return u.tokens, nil
}

func (u *InMemoryStore) RetrieveToken(token_id uuid.UUID) (*Token, error) {
	for _, token := range u.tokens {
		if token.ID == token_id {
			return &token, nil
		}
	}
	return nil, ErrNotFound
}

func (u *InMemoryStore) UpdateToken(token *Token) error {
	for idx, t := range u.tokens {
		if t.ID == token.ID {
			u.tokens[idx] = *token
			return nil
		}
	}
	return ErrNotFound
}

func (u *InMemoryStore) Persist() error {
	data := struct {
		Users  []User  `json:"users"`
		Tokens []Token `json:"tokens"`
	}{
		Users:  u.users,
		Tokens: u.tokens,
	}

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(FILE_NAME, jsonData, 0644)
}

func (u *InMemoryStore) Load() error {
	jsonData, err := os.ReadFile(FILE_NAME)
	if err != nil {
		return err
	}

	var data struct {
		Users  []User  `json:"users"`
		Tokens []Token `json:"tokens"`
	}

	err = json.Unmarshal(jsonData, &data)
	if err != nil {
		return err
	}

	u.users = data.Users
	u.tokens = data.Tokens

	return nil
}
