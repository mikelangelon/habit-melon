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
	postHabit(t, s, newHabit)
	habit := getHabit(t, s, newHabit.Description)
	assert.Equal(t, newHabit, habit)
}

func TestPutHabit(t *testing.T) {
	s := api.NewServer()
	newHabit := api.Habit{
		Description: "Test",
		Days:        []civil.Date{{Day: 8, Month: time.April, Year: 2023}},
	}
	postHabit(t, s, newHabit)
	changedHabit := api.Habit{
		Description: "Test",
		Days: []civil.Date{
			{Day: 8, Month: time.April, Year: 2023},
			{Day: 9, Month: time.April, Year: 2023},
		},
	}
	putHabit(t, s, changedHabit)
	got := getHabit(t, s, changedHabit.Description)
	assert.Equal(t, changedHabit, got)
}

func TestPutHabitNotFound(t *testing.T) {
	s := api.NewServer()
	h := api.Habit{
		Description: "Test",
		Days: []civil.Date{
			{Day: 8, Month: time.April, Year: 2023},
			{Day: 9, Month: time.April, Year: 2023},
		},
	}
	request, err := json.Marshal(h)
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest("PUT", fmt.Sprintf("/v1/habit/%s", h.Description), bytes.NewReader(request))
	if err != nil {
		t.Fatal(err)
	}
	recorder := httptest.NewRecorder()
	s.ServeHTTP(recorder, req)
	statusCode := http.StatusNotFound
	if recorder.Result().StatusCode != statusCode {
		t.Errorf("an unexpected result updating habit: got %v want %v", recorder.Result().StatusCode, statusCode)
	}
}
func TestDeleteHabit(t *testing.T) {
	s := api.NewServer()
	newHabit := api.Habit{
		Description: "Test",
		Days:        []civil.Date{{Day: 8, Month: time.April, Year: 2023}},
	}
	postHabit(t, s, newHabit)

	req, err := http.NewRequest("DELETE", fmt.Sprintf("/v1/habit/%s", newHabit.Description), http.NoBody)
	if err != nil {
		t.Fatal(err)
	}
	recorder := httptest.NewRecorder()
	s.ServeHTTP(recorder, req)

	statusCode := 200
	if recorder.Result().StatusCode != statusCode {
		t.Errorf("TestDeleteHabit() test returned an unexpected result: got %v want %v", recorder.Result().StatusCode, statusCode)
	}
}
func TestDeleteHabitNotFound(t *testing.T) {
	s := api.NewServer()

	req, err := http.NewRequest("DELETE", "/v1/habit/something", http.NoBody)
	if err != nil {
		t.Fatal(err)
	}
	recorder := httptest.NewRecorder()
	s.ServeHTTP(recorder, req)

	statusCode := http.StatusNotFound
	if recorder.Result().StatusCode != statusCode {
		t.Errorf("TestDeleteHabit() test returned an unexpected result: got %v want %v", recorder.Result().StatusCode, statusCode)
	}
}

// Helpers
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
func postHabit(t *testing.T, s *api.Server, h api.Habit) {
	request, err := json.Marshal(h)
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest("POST", "/v1/habit", bytes.NewReader(request))
	if err != nil {
		t.Fatal(err)
	}
	recorder := httptest.NewRecorder()
	s.ServeHTTP(recorder, req)
	statusCode := http.StatusCreated
	if recorder.Result().StatusCode != statusCode {
		t.Errorf("an unexpected result creating habit: got %v want %v", recorder.Result().StatusCode, statusCode)
	}
}
func putHabit(t *testing.T, s *api.Server, h api.Habit) {
	request, err := json.Marshal(h)
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest("PUT", fmt.Sprintf("/v1/habit/%s", h.Description), bytes.NewReader(request))
	if err != nil {
		t.Fatal(err)
	}
	recorder := httptest.NewRecorder()
	s.ServeHTTP(recorder, req)
	statusCode := http.StatusOK
	if recorder.Result().StatusCode != statusCode {
		t.Errorf("an unexpected result updating habit: got %v want %v", recorder.Result().StatusCode, statusCode)
	}
}
