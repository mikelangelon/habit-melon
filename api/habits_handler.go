package api

import (
	"cloud.google.com/go/civil"
	"net/http"
)

type HabitsResponse struct {
	Habits []Habit
}

type Habit struct {
	Description string
	Days        []civil.Date
}

func (hr HabitsResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
