package main

import (
	"fmt"
	"microservices-project/symbol-service/symbols"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func mustConnectionString() string {
	pgUser := os.Getenv("POSTGRES_USER")
	pgPassword := os.Getenv("POSTGRES_PASSWORD")
	pgDatabase := os.Getenv("POSTGRES_DB")
	pgHost := os.Getenv("POSTGRES_HOST")
	// TODO: get rid of sslmode disable in production
	return fmt.Sprintf("user=%v password=%v dbname=%v host=%v sslmode=disable", pgUser, pgPassword, pgDatabase, pgHost)
}

func main() {
	r := mux.NewRouter()
	connStr := mustConnectionString()
	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		panic(err)
	}
	symbols.RegisterRoutes(r, db)
	http.ListenAndServe(":5050", r)
}
