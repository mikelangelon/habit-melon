package main

import (
	"cloud.google.com/go/civil"
	"context"
	"github.com/mikelangelon/habit-melon/api"
	"github.com/mikelangelon/habit-melon/internal/core"
	"github.com/mikelangelon/habit-melon/internal/postgres"
	"time"
)

func main() {
	db, err := postgres.NewPostgresDB("postgresql://db-user:db-pass@localhost:5432/db-name?sslmode=disable")
	if err != nil {
		panic(err)
	}

	repo := postgres.New(db)
	_, err = repo.CreateHabit(context.TODO(), core.Habit{
		Description: "Test8",
		Days: []civil.Date{
			{Day: 1, Month: time.April, Year: 2023},
		},
	})
	if err != nil {
		panic(err)
	}
	server := api.NewServer()
	server.Start()
}
