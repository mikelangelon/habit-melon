package api

import (
	"cloud.google.com/go/civil"
)

type Habit struct {
	Description string
	Days        []civil.Date
}
