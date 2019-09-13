package main

import (
	"fmt"
	"microservices-project/quote-service/quotes"
	"os"

	"github.com/adeaver/microservices-project/util/httpservice"
	"github.com/jmoiron/sqlx"
	"github.com/streadway/amqp"

	_ "github.com/lib/pq"
)

func mustPostgresConnectionString() string {
	pgUser := os.Getenv("POSTGRES_USER")
	pgPassword := os.Getenv("POSTGRES_PASSWORD")
	pgDatabase := os.Getenv("POSTGRES_DB")
	pgHost := os.Getenv("POSTGRES_HOST")
	// TODO: get rid of sslmode disable in production
	return fmt.Sprintf("user=%v password=%v dbname=%v host=%v sslmode=disable", pgUser, pgPassword, pgDatabase, pgHost)
}

func mustRabbitMQConnectionString() string {
	amqpUser := os.Getenv("RABBITMQ_USER")
	amqpPassword := os.Getenv("RABBITMQ_PASSWORD")
	amqpHost := os.Getenv("RABBITMQ_HOST")
	amqpPort := os.Getenv("RABBITMQ_PORT")
	return fmt.Sprintf("amqp://%s:%s@%s:%s/", amqpUser, amqpPassword, amqpHost, amqpPort)
}

func main() {
	pgConnStr := mustPostgresConnectionString()
	db, err := sqlx.Connect("postgres", pgConnStr)
	if err != nil {
		panic(err)
	}
	brokerConnStr := mustRabbitMQConnectionString()
	brokerConn, err := amqp.Dial(brokerConnStr)
	if err != nil {
		panic(err)
	}
	ch, err := brokerConn.Channel()
	if err != nil {
		panic(err)
	}
	s := httpservice.Service{
		Routes:     quotes.MakeRouteDefinitions(db, ch),
		ListenPort: ":3050",
	}
	s.Start()
}
