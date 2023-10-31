package main

import (
	"github.com/mikelangelon/habit-melon/api"
)

func main() {
	server := api.NewServer()
	server.Start()
}
