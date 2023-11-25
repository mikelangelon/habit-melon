package main

import "fmt"

var (
	Version = "development"
	DBUser  = "db-user"
	DBPass  = "db-pass"
	DBHost  = "localhost:5432"
	DBName  = "db-name"
)

func DBUrl() string {
	return fmt.Sprintf("postgresql://%s:%s@%s/%s?sslmode=disable", DBUser, DBPass, DBHost, DBName)
}
