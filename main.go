package main

import (
	"github.com/mikelangelon/habit-melon/api"
	"github.com/mikelangelon/habit-melon/internal/app"
	"github.com/mikelangelon/habit-melon/internal/postgres"
)

func main() {
	db, err := postgres.NewPostgresDB("postgresql://db-user:db-pass@localhost:5432/db-name?sslmode=disable")
	if err != nil {
		panic(err)
	}

	repo := postgres.New(db)
	service := app.NewHabitService(repo)
	server := api.NewServer(service)
	server.Start()
}
