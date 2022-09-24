package main

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func apiRequest(t *testing.T, f http.HandlerFunc, args map[string]string) (int, []byte) {
	w := httptest.NewRecorder()

	form := url.Values{}
	for key, arg := range args {
		form.Set(key, arg)
	}

	req, err := http.NewRequest("POST", "/", strings.NewReader(form.Encode()))
	require.NoError(t, err)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	f(w, req)

	return w.Code, w.Body.Bytes()
}

// Integration test? What is a strategy?
func TestHandleEmailRegister(t *testing.T) {
	setupRedis(t)

	// Register
	status, respRegBytes := apiRequest(t, handleEmailRegister, map[string]string{
		"email":    "my@email.com",
		"password": "1234",
	})

	respReg := struct {
		Token string
		User  struct {
			ID string
		}
	}{}
	err := json.Unmarshal(respRegBytes, &respReg)
	require.NoError(t, err)

	require.Equal(t, 200, status)
	require.NotEmpty(t, respReg.Token)
	require.NotEmpty(t, respReg.User.ID)

	// Login
	status, body := apiRequest(t, handleAuthEmail, map[string]string{
		"email":    "my@email.com",
		"password": "1234",
	})
	require.Equal(t, 200, status)

	respLogin := struct {
		Token string
		User  struct {
			ID string
		}
	}{}
	err = json.Unmarshal(body, &respLogin)
	require.NoError(t, err)
	require.NotEmpty(t, respLogin.Token)
	require.Equal(t, respReg.User.ID, respLogin.User.ID)
}
