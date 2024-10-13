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

func TestAuth(t *testing.T) {
	store, err := NewInMemoryStore()
	require.NoError(t, err)
	logger := log.New(os.Stdout, "test: ", log.LstdFlags)
	app := application{
		store:  store,
		logger: logger,
	}
	mux := app.mountRoutes()
	t.Run("test login for non existent user", func(t *testing.T) {
		var login_buffer bytes.Buffer
		loginRequest := loginRequest{
			Email:    "test@user.com",
			Password: "password",
		}
		err = json.NewEncoder(&login_buffer).Encode(loginRequest)
		createRequest, err := http.NewRequest(http.MethodPost, "/api/v1/auth/login", &login_buffer)
		require.NoError(t, err)

		rr := httptest.NewRecorder()

		mux.ServeHTTP(rr, createRequest)
		require.Equal(t, http.StatusBadRequest, rr.Code)
		var resp ErrorMessage
		err = json.NewDecoder(rr.Body).Decode(&resp)
		require.Equal(t, "invalid email or password provided", resp.Error)
	})
	t.Run("test login", func(t *testing.T) {
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
		var login_buffer bytes.Buffer
		loginRequest := loginRequest{
			Email:    "test@user.com",
			Password: "password",
		}
		err = json.NewEncoder(&login_buffer).Encode(loginRequest)
		createRequest, err = http.NewRequest(http.MethodPost, "/api/v1/auth/login", &login_buffer)
		require.NoError(t, err)

		rr = httptest.NewRecorder()

		mux.ServeHTTP(rr, createRequest)
		require.Equal(t, http.StatusOK, rr.Code)
		var resp Token
		err = json.NewDecoder(rr.Body).Decode(&resp)
		require.NotEmpty(t, resp.Token)
		require.NotEmpty(t, resp.CreatedAt)
		require.NotEmpty(t, resp.Expiry)
		require.Equal(t, resp.Revoked, false)
		require.NotEmpty(t, resp.UserId)
	})
	t.Run("test login with incorrect password", func(t *testing.T) {
		var login_buffer bytes.Buffer
		loginRequest := loginRequest{
			Email:    "test@user.com",
			Password: "password123",
		}
		err = json.NewEncoder(&login_buffer).Encode(loginRequest)
		createRequest, err := http.NewRequest(http.MethodPost, "/api/v1/auth/login", &login_buffer)
		require.NoError(t, err)

		rr := httptest.NewRecorder()

		mux.ServeHTTP(rr, createRequest)
		require.Equal(t, http.StatusBadRequest, rr.Code)
		var resp ErrorMessage
		err = json.NewDecoder(rr.Body).Decode(&resp)
		require.Equal(t, "invalid email or password provided", resp.Error)
	})
}
