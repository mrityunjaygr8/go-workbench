package main

const FILE_NAME = "store.json"

func NewInMemoryStore() (store, error) {
	users := make([]User, 0)
	tokens := make([]Token, 0)

	return &InMemoryStore{
		users:  users,
		tokens: tokens,
	}, nil
}

type InMemoryStore struct {
	users  []User
	tokens []Token
}

func (u *InMemoryStore) ListUsers() ([]User, error) {
	return u.users, nil
}

func (u *InMemoryStore) InsertUser(user User) error {
	u.users = append(u.users, user)
	return nil
}

func (u *InMemoryStore) GetUserByEmail(email string) (*User, error) {
	for _, user := range u.users {
		if user.Email == email {
			return &user, nil
		}
	}
	return nil, ErrUserNotFound
}

func (u *InMemoryStore) InsertToken(token Token) error {
	u.tokens = append(u.tokens, token)
	return nil

}

func (u *InMemoryStore) ListTokens() ([]Token, error) {
	return u.tokens, nil
}
