package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAPIClient_Authenticate(t *testing.T) {
	// Создаем имитированный сервер
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/authentication", r.URL.Path)

		var data map[string]string
		err := json.NewDecoder(r.Body).Decode(&data)
		assert.NoError(t, err)
		assert.Equal(t, "user", data["login"])
		assert.Equal(t, "password", data["password"])

		w.Header().Set("Set-Cookie", "user_info=test_token")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "success"})
	}))
	defer ts.Close()

	client := NewAPIClient(ts.URL)
	token, responseBody, err := client.Authenticate("user", "password")
	assert.NoError(t, err)
	assert.Equal(t, "test_token", token)

	var response map[string]string
	err = json.Unmarshal(responseBody, &response)
	assert.NoError(t, err)
	assert.Equal(t, "success", response["status"])
}

func TestAPIClient_Registration(t *testing.T) {
	// Создаем имитированный сервер
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/registration", r.URL.Path)

		var data map[string]string
		err := json.NewDecoder(r.Body).Decode(&data)
		assert.NoError(t, err)
		assert.Equal(t, "user", data["login"])
		assert.Equal(t, "password", data["password"])

		w.Header().Set("Set-Cookie", "user_info=test_token")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "registered"})
	}))
	defer ts.Close()

	client := NewAPIClient(ts.URL)
	token, responseBody, err := client.Registration("user", "password")
	assert.NoError(t, err)
	assert.Equal(t, "test_token", token)

	var response map[string]string
	err = json.Unmarshal(responseBody, &response)
	assert.NoError(t, err)
	assert.Equal(t, "registered", response["status"])
}

func TestAPIClient_Get(t *testing.T) {
	// Создаем имитированный сервер
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/test-endpoint", r.URL.Path)
		assert.Equal(t, "user_info=test_token", r.Header.Get("Cookie"))

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "success"})
	}))
	defer ts.Close()

	client := NewAPIClient(ts.URL)
	client.SetToken("test_token")
	responseBody, err := client.Get("test-endpoint", nil)
	assert.NoError(t, err)

	var response map[string]string
	err = json.Unmarshal(responseBody, &response)
	assert.NoError(t, err)
	assert.Equal(t, "success", response["message"])
}

func TestAPIClient_Post(t *testing.T) {
	// Создаем имитированный сервер
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/test-endpoint", r.URL.Path)
		assert.Equal(t, "user_info=test_token", r.Header.Get("Cookie"))

		var data map[string]string
		err := json.NewDecoder(r.Body).Decode(&data)
		assert.NoError(t, err)
		assert.Equal(t, "testData", data["data"])

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "success"})
	}))
	defer ts.Close()

	client := NewAPIClient(ts.URL)
	client.SetToken("test_token")
	responseBody, err := client.Post("test-endpoint", map[string]string{"data": "testData"}, nil)
	assert.NoError(t, err)

	var response map[string]string
	err = json.Unmarshal(responseBody, &response)
	assert.NoError(t, err)
	assert.Equal(t, "success", response["message"])
}

func TestAPIClient_Ping(t *testing.T) {
	// Создаем имитированный сервер
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/ping", r.URL.Path)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	client := NewAPIClient(ts.URL)
	err := client.Ping()
	assert.NoError(t, err)
}
