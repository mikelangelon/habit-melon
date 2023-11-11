package core

import "cloud.google.com/go/civil"

type Habit struct {
	ID          *int64
	Description string
	Days        []civil.Date
}
