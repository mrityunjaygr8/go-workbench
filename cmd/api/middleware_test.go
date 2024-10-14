package main

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTracer(t *testing.T) {
	mux := http.NewServeMux()
	testHandler := func(w http.ResponseWriter, r *http.Request) {
	}

	reqFn := func() string {
		return "asdf"
	}
	mux.HandleFunc("/", testHandler)
	testRequest, err := http.NewRequest(http.MethodGet, "/", nil)
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	tracing(reqFn)(mux).ServeHTTP(rr, testRequest)

	require.NotEmpty(t, rr.Header().Get("X-Request-Id"))
	require.Equal(t, "asdf", rr.Header().Get("X-Request-Id"))
}

func TestLogging(t *testing.T) {
	mux := http.NewServeMux()
	testHandler := func(w http.ResponseWriter, r *http.Request) {

	}
	mux.HandleFunc("/", testHandler)
	testRequest, err := http.NewRequest(http.MethodGet, "/", nil)
	require.NoError(t, err)

	var mockedLogBuf bytes.Buffer
	mockedLog := log.New(&mockedLogBuf, "woo ", log.LstdFlags)

	rr := httptest.NewRecorder()
	logging(mockedLog)(mux).ServeHTTP(rr, testRequest)

	require.NotEmpty(t, mockedLogBuf)
	t.Log(string(mockedLogBuf.Bytes()))
}
