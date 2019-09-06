package main

import (
	"fmt"
	"microservices-project/symbol-service/symbols"
	"os"

	"github.com/adeaver/microservices-project/util/httpservice"
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
	connStr := mustConnectionString()
	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		panic(err)
	}
	s := httpservice.Service{
		Routes:     symbols.MakeRouteDefinitions(db),
		ListenPort: ":5050",
	}
	s.Start()
}
