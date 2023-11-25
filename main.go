package main

import (
	"github.com/mikelangelon/habit-melon/api"
	"github.com/mikelangelon/habit-melon/internal/app"
	"github.com/mikelangelon/habit-melon/internal/postgres"
	"log/slog"
)

func main() {
	db, err := postgres.NewPostgresDB(DBUrl())
	if err != nil {
		panic(err)
	}

	slog.Info("Starting server", "version", Version)
	repo := postgres.New(db)
	service := app.NewHabitService(repo)
	server := api.NewServer(service)
	server.Start()
}
