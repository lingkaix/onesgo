package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestServer(t *testing.T) {
	r := setupRouter()
	var token string
	var userID uint64

	t.Run("register new user", func(t *testing.T) {
		regData := gin.H{
			"email":      "test@example.com",
			"password":   "testpassword",
			"first_name": "Lingkai",
			"last_name":  "Xu",
		}
		regBytes, err := json.Marshal(regData)
		assert.NoError(t, err)

		req, err := http.NewRequest("POST", "/register", bytes.NewReader(regBytes))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
		var resp User
		err = json.NewDecoder(strings.NewReader(w.Body.String())).Decode(&resp)
		assert.NoError(t, err)
		userID = resp.ID
	})
	t.Run("user login", func(t *testing.T) {
		loginData := gin.H{
			"email":    "test@example.com",
			"password": "testpassword",
		}
		loginBytes, err := json.Marshal(loginData)
		assert.NoError(t, err)

		req, err := http.NewRequest("POST", "/login", bytes.NewReader(loginBytes))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
		var resp struct {
			Token string `json:"token"`
		}
		err = json.NewDecoder(strings.NewReader(w.Body.String())).Decode(&resp)
		assert.NoError(t, err)
		token = resp.Token
	})

	t.Run("get user info", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/users/"+strconv.FormatUint(userID, 10), nil)
		assert.NoError(t, err)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
		var resp User
		err = json.NewDecoder(strings.NewReader(w.Body.String())).Decode(&resp)
		assert.NoError(t, err)
		assert.Equal(t, userID, resp.ID)
	})

	t.Run("update user", func(t *testing.T) {
		phone := "1234567890"
		fName := "Simon"
		uData := gin.H{
			"first_name": fName,
			"phone":      phone,
		}
		uBytes, err := json.Marshal(uData)
		assert.NoError(t, err)

		req, err := http.NewRequest("POST", "/users/"+strconv.FormatUint(userID, 10), bytes.NewReader(uBytes))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
		var resp User
		err = json.NewDecoder(strings.NewReader(w.Body.String())).Decode(&resp)
		assert.NoError(t, err)
		assert.Equal(t, userID, resp.ID)
		assert.Equal(t, phone, resp.Phone)
		assert.Equal(t, fName, resp.FirstName)
	})

	t.Run("get user info", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/users/"+strconv.FormatUint(userID, 10), nil)
		assert.NoError(t, err)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
		var resp User
		err = json.NewDecoder(strings.NewReader(w.Body.String())).Decode(&resp)
		assert.NoError(t, err)
		assert.Equal(t, userID, resp.ID)
	})

	t.Run("delete user", func(t *testing.T) {
		req, err := http.NewRequest("DELETE", "/users/"+strconv.FormatUint(userID, 10), nil)
		assert.NoError(t, err)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
		var resp struct {
			Status string `json:"status"`
		}
		err = json.NewDecoder(strings.NewReader(w.Body.String())).Decode(&resp)
		assert.NoError(t, err)
		assert.Equal(t, "deleted", resp.Status)
	})
}
