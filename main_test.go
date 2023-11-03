package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mikelangelon/habit-melon/api"
)

func TestHelloWorld(t *testing.T) {
	// Set up a new request.
	req, err := http.NewRequest("GET", "/", http.NoBody)
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	api.NewServer().ServeHTTP(recorder, req)

	statusCode := 200
	if recorder.Result().StatusCode != statusCode {
		t.Errorf("TestInfoRequest() test returned an unexpected result: got %v want %v", recorder.Result().StatusCode, statusCode)
	}
	var response string
	body := recorder.Body.String()
	if err != nil {
		t.Fatalf("reading response body: %v", err)
	}
	if body != "Hello world" {
		t.Fatalf("bad response message: %s", response)
	}
}

func TestHealth(t *testing.T) {
	// Set up a new request.
	req, err := http.NewRequest("GET", "/health", http.NoBody)
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	api.NewServer().ServeHTTP(recorder, req)

	statusCode := 200
	if recorder.Result().StatusCode != statusCode {
		t.Errorf("TestInfoRequest() test returned an unexpected result: got %v want %v", recorder.Result().StatusCode, statusCode)
	}
	var response api.HealthResponse
	body := recorder.Body.Bytes()
	if err != nil {
		t.Fatalf("reading response body: %v", err)
	}
	if err := json.Unmarshal(body, &response); err != nil {
		t.Fatalf("parsing json response: %v", err)
	}
	if response.OK != true {
		t.Fatalf("bad response message: %v", response)
	}
}

func TestGetHabit(t *testing.T) {
	// Set up a new request.
	req, err := http.NewRequest("GET", "/v1/habit", http.NoBody)
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	api.NewServer().ServeHTTP(recorder, req)

	statusCode := 200
	if recorder.Result().StatusCode != statusCode {
		t.Errorf("TestInfoRequest() test returned an unexpected result: got %v want %v", recorder.Result().StatusCode, statusCode)
	}
	var response string
	body := recorder.Body.String()
	if err != nil {
		t.Fatalf("reading response body: %v", err)
	}
	if body != "Hello world" {
		t.Fatalf("bad response message: %s", response)
	}
}
