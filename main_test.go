package main

import (
	"bytes"
	"cloud.google.com/go/civil"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

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
		t.Errorf("TestHelloWorld() test returned an unexpected result: got %v want %v", recorder.Result().StatusCode, statusCode)
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
		t.Errorf("TestHealth() test returned an unexpected result: got %v want %v", recorder.Result().StatusCode, statusCode)
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
	s := api.NewServer()
	// Set up a new request.
	got := getHabit(t, s, "Study Dutch")
	assert.Equal(t, api.Habit{
		Description: "Study Dutch",
		Days: []civil.Date{
			{Day: 1, Month: time.April, Year: 2023},
			{Day: 2, Month: time.April, Year: 2023},
			{Day: 3, Month: time.April, Year: 2023},
		},
	}, got)
}

func TestGetHabitNoFound(t *testing.T) {
	// Set up a new request.
	req, err := http.NewRequest("GET", "/v1/habit/NotFound", http.NoBody)
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	api.NewServer().ServeHTTP(recorder, req)

	statusCode := 404
	if recorder.Result().StatusCode != statusCode {
		t.Errorf("TestGetHabit() test returned an unexpected result: got %v want %v", recorder.Result().StatusCode, statusCode)
	}
}

func TestPostHabit(t *testing.T) {
	s := api.NewServer()
	newHabit := api.Habit{
		Description: "Test",
		Days:        []civil.Date{{Day: 8, Month: time.April, Year: 2023}},
	}
	// Set up a new request.
	request, err := json.Marshal(newHabit)
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest("POST", "/v1/habit", bytes.NewReader(request))
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	s.ServeHTTP(recorder, req)

	statusCode := 201
	if recorder.Result().StatusCode != statusCode {
		t.Errorf("TestPostHabit() test returned an unexpected result: got %v want %v", recorder.Result().StatusCode, statusCode)
	}
	habit := getHabit(t, s, newHabit.Description)
	assert.Equal(t, newHabit, habit)
}

func getHabit(t *testing.T, s *api.Server, id string) api.Habit {
	req, err := http.NewRequest("GET", fmt.Sprintf("/v1/habit/%s", id), http.NoBody)
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()
	s.ServeHTTP(recorder, req)

	statusCode := 200
	if recorder.Result().StatusCode != statusCode {
		t.Errorf("TestGetHabit() test returned an unexpected result: got %v want %v", recorder.Result().StatusCode, statusCode)
	}
	var response api.Habit
	body := recorder.Body.Bytes()
	err = json.Unmarshal(body, &response)
	if err != nil {
		t.Fatalf("reading response body: %v", err)
	}
	return response
}
