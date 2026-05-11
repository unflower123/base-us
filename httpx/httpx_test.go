package httpx

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// TestGetRequest test GET request
func TestGetRequest(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "success", "code": 200}`))
	}))
	defer ts.Close()

	client := NewHttpClient(10*time.Second, 3, 1*time.Second)

	var response struct {
		Message string `json:"message"`
		Code    int    `json:"code"`
	}

	err := client.Get(context.Background(), ts.URL, nil, &response)
	if err != nil {
		t.Errorf("GET request failed: %v", err)
	}

	if response.Message != "success" || response.Code != 200 {
		t.Errorf("Unexpected response: %+v", response)
	}
}

// TestPostRequest test POST request
func TestPostRequest(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		var body struct {
			Name string `json:"name"`
			Age  int    `json:"age"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Errorf("Failed to decode request body: %v", err)
		}

		if body.Name != "John" || body.Age != 30 {
			t.Errorf("Unexpected request body: %+v", body)
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "success", "code": 200}`))
	}))
	defer ts.Close()

	client := NewHttpClient(10*time.Second, 3, 1*time.Second)

	requestBody := struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}{
		Name: "John",
		Age:  30,
	}

	var response struct {
		Message string `json:"message"`
		Code    int    `json:"code"`
	}

	err := client.Post(context.Background(), ts.URL, nil, requestBody, &response)
	if err != nil {
		t.Errorf("POST request failed: %v", err)
	}

	if response.Message != "success" || response.Code != 200 {
		t.Errorf("Unexpected response: %+v", response)
	}
}

// TestRetryMechanism test retry mechanism
func TestRetryMechanism(t *testing.T) {
	attempt := 0

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempt++
		if attempt < 3 { // Return 500 error in the first two attempts
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "success", "code": 200}`))
	}))
	defer ts.Close()

	client := NewHttpClient(10*time.Second, 3, 1*time.Second)

	var response struct {
		Message string `json:"message"`
		Code    int    `json:"code"`
	}

	err := client.Get(context.Background(), ts.URL, nil, &response)
	if err != nil {
		t.Errorf("GET request failed: %v", err)
	}

	if response.Message != "success" || response.Code != 200 {
		t.Errorf("Unexpected response: %+v", response)
	}

	if attempt != 3 {
		t.Errorf("Expected 3 attempts, got %d", attempt)
	}
}

// TestContextCancel test context cancel
func TestContextCancel(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second) // 模拟延迟
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	client := NewHttpClient(10*time.Second, 3, 1*time.Second)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	var response struct {
		Message string `json:"message"`
		Code    int    `json:"code"`
	}

	err := client.Get(ctx, ts.URL, nil, &response)
	if err == nil {
		t.Error("Expected context canceled error, got nil")
	}
}
