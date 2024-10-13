package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUsers(t *testing.T) {
	store, err := NewInMemoryStore()
	require.NoError(t, err)
	logger := log.New(os.Stdout, "test: ", log.LstdFlags)
	app := application{
		store:  store,
		logger: logger,
	}
	mux := app.mountRoutes()
	t.Run("test user creation", func(t *testing.T) {
		var b bytes.Buffer
		userCreationRequest := createUserRequest{
			Email:     "test@user.com",
			Password:  "password",
			FirstName: "FirstName",
			LastName:  "LastName",
		}
		err := json.NewEncoder(&b).Encode(userCreationRequest)
		createRequest, err := http.NewRequest(http.MethodPost, "/api/v1/users", &b)
		require.NoError(t, err)

		rr := httptest.NewRecorder()

		mux.ServeHTTP(rr, createRequest)
		require.Equal(t, http.StatusCreated, rr.Code)

	})
	t.Run("test user listing", func(t *testing.T) {
		createRequest, err := http.NewRequest(http.MethodGet, "/api/v1/users", nil)
		require.NoError(t, err)

		rr := httptest.NewRecorder()

		mux.ServeHTTP(rr, createRequest)
		require.Equal(t, http.StatusOK, rr.Code)
		var resp []User
		err = json.NewDecoder(rr.Body).Decode(&resp)
		require.Equal(t, 1, len(resp))
		require.Equal(t, "test@user.com", resp[0].Email)
		require.Equal(t, "FirstName", resp[0].FirstName)
		require.Equal(t, "LastName", resp[0].LastName)
		require.NotEmpty(t, resp[0].ID)
	})
	t.Run("test same email user creation", func(t *testing.T) {
		var b bytes.Buffer
		userCreationRequest := createUserRequest{
			Email:     "test@user.com",
			Password:  "password",
			FirstName: "FirstName",
			LastName:  "LastName",
		}
		err := json.NewEncoder(&b).Encode(userCreationRequest)
		createRequest, err := http.NewRequest(http.MethodPost, "/api/v1/users", &b)
		require.NoError(t, err)

		rr := httptest.NewRecorder()

		mux.ServeHTTP(rr, createRequest)

		require.Equal(t, http.StatusBadRequest, rr.Code)
		var resp ErrorMessage
		err = json.NewDecoder(rr.Body).Decode(&resp)
		require.Equal(t, "user already exists", resp.Error)
	})
}
