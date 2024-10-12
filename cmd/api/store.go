package main

import "github.com/google/uuid"

type store struct {
	Users interface {
		ListUsers() ([]User, error)
		CreateUser(user User) error
		GetUserByEmail(email string) (*User, error)
	}
}

func NewInMemoryStore() store {
	dict := make(map[string]any)
	dict["users"] = make([]User, 0)
	return store{
		Users: &InMemoryStoreUsers{
			storage: dict["users"].([]User),
		},
	}
}

type User struct {
	Email     string    `json:"email"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Password  string    `json:"-"`
	ID        uuid.UUID `json:"id"`
}

type InMemoryStore struct {
	storage map[string]any
}

type InMemoryStoreUsers struct {
	storage []User
}

func (u *InMemoryStoreUsers) ListUsers() ([]User, error) {
	return u.storage, nil
}

func (u *InMemoryStoreUsers) CreateUser(user User) error {
	uuid, err := uuid.NewV7()
	if err != nil {
		return err
	}

	user.ID = uuid
	u.storage = append(u.storage, user)
	return nil
}

func (u *InMemoryStoreUsers) GetUserByEmail(email string) (*User, error) {
	for _, user := range u.storage {
		if user.Email == email {
			return &user, nil
		}
	}
	return nil, nil
}
