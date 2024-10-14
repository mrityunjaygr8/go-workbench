package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/google/uuid"
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
		var loginBuffer bytes.Buffer
		loginRequestPayload := loginRequest{
			Email:    "test@user.com",
			Password: "password",
		}
		err = json.NewEncoder(&loginBuffer).Encode(loginRequestPayload)
		loginRequest, err := http.NewRequest(http.MethodPost, "/api/v1/auth/login", &loginBuffer)
		require.NoError(t, err)

		rr := httptest.NewRecorder()

		mux.ServeHTTP(rr, loginRequest)
		require.Equal(t, http.StatusBadRequest, rr.Code)
		var resp ErrorMessage
		err = json.NewDecoder(rr.Body).Decode(&resp)
		require.Equal(t, "invalid email or password provided", resp.Error)
	})
	t.Run("test login", func(t *testing.T) {
		var b bytes.Buffer
		userCreationRequestPayload := createUserRequest{
			Email:     "test@user.com",
			Password:  "password",
			FirstName: "FirstName",
			LastName:  "LastName",
		}
		err := json.NewEncoder(&b).Encode(userCreationRequestPayload)
		createRequest, err := http.NewRequest(http.MethodPost, "/api/v1/users", &b)
		require.NoError(t, err)

		rr := httptest.NewRecorder()

		mux.ServeHTTP(rr, createRequest)
		require.Equal(t, http.StatusCreated, rr.Code)
		var login_buffer bytes.Buffer
		loginRequestPayload := loginRequest{
			Email:    "test@user.com",
			Password: "password",
		}
		err = json.NewEncoder(&login_buffer).Encode(loginRequestPayload)
		loginRequest, err := http.NewRequest(http.MethodPost, "/api/v1/auth/login", &login_buffer)
		require.NoError(t, err)

		rr = httptest.NewRecorder()

		mux.ServeHTTP(rr, loginRequest)
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
		var loginBuffer bytes.Buffer
		loginRequestPayload := loginRequest{
			Email:    "test@user.com",
			Password: "password123",
		}
		err = json.NewEncoder(&loginBuffer).Encode(loginRequestPayload)
		loginRequest, err := http.NewRequest(http.MethodPost, "/api/v1/auth/login", &loginBuffer)
		require.NoError(t, err)

		rr := httptest.NewRecorder()

		mux.ServeHTTP(rr, loginRequest)
		require.Equal(t, http.StatusBadRequest, rr.Code)
		var resp ErrorMessage
		err = json.NewDecoder(rr.Body).Decode(&resp)
		require.Equal(t, "invalid email or password provided", resp.Error)
	})
	t.Run("test token revoke", func(t *testing.T) {
		var loginBuffer bytes.Buffer
		loginRequestPayload := loginRequest{
			Email:    "test@user.com",
			Password: "password",
		}
		err = json.NewEncoder(&loginBuffer).Encode(loginRequestPayload)
		loginRequest, err := http.NewRequest(http.MethodPost, "/api/v1/auth/login", &loginBuffer)
		require.NoError(t, err)

		rr := httptest.NewRecorder()

		mux.ServeHTTP(rr, loginRequest)
		require.Equal(t, http.StatusOK, rr.Code)
		var resp Token
		err = json.NewDecoder(rr.Body).Decode(&resp)
		var revokeBuffer bytes.Buffer
		tokenRevokeRequest := tokenRevokeRequest{
			TokenID: resp.ID.String(),
			// TokenID: "",
		}
		err = json.NewEncoder(&revokeBuffer).Encode(tokenRevokeRequest)
		revokeRequest, err := http.NewRequest(http.MethodPost, "/api/v1/auth/revoke", &revokeBuffer)
		require.NoError(t, err)

		rr = httptest.NewRecorder()

		mux.ServeHTTP(rr, revokeRequest)

		require.Equal(t, http.StatusOK, rr.Code)

		token, err := store.RetrieveToken(resp.ID)
		require.NoError(t, err)
		require.Equal(t, true, token.Revoked)
	})
	t.Run("test token revoke invalid uuid", func(t *testing.T) {
		var revokeBuffer bytes.Buffer
		tokenRevokeRequest := tokenRevokeRequest{
			TokenID: "",
		}
		err = json.NewEncoder(&revokeBuffer).Encode(tokenRevokeRequest)
		revokeRequest, err := http.NewRequest(http.MethodPost, "/api/v1/auth/revoke", &revokeBuffer)
		require.NoError(t, err)

		rr := httptest.NewRecorder()

		mux.ServeHTTP(rr, revokeRequest)

		require.Equal(t, http.StatusBadRequest, rr.Code)
		var resp ErrorMessage
		err = json.NewDecoder(rr.Body).Decode(&resp)
		require.Equal(t, "not a valid uuid", resp.Error)
	})
	t.Run("test token revoke unknown uuid", func(t *testing.T) {
		var revokeBuffer bytes.Buffer
		tokenRevokeRequest := tokenRevokeRequest{
			TokenID: uuid.New().String(),
		}
		err = json.NewEncoder(&revokeBuffer).Encode(tokenRevokeRequest)
		revokeRequest, err := http.NewRequest(http.MethodPost, "/api/v1/auth/revoke", &revokeBuffer)
		require.NoError(t, err)

		rr := httptest.NewRecorder()

		mux.ServeHTTP(rr, revokeRequest)

		require.Equal(t, http.StatusBadRequest, rr.Code)
		var resp ErrorMessage
		err = json.NewDecoder(rr.Body).Decode(&resp)
		require.Equal(t, "not found", resp.Error)
	})
}
