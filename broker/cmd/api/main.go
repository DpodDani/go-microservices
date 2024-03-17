package main

import (
	"fmt"
	"log"
	"net/http"

	rmq "github.com/DpodDani/go-microservices-toolbox/rmq"
	amqp "github.com/rabbitmq/amqp091-go"
)

const webPort = "8080"

type Config struct {
	rmqConn *amqp.Connection
}

func main() {
	conn, err := rmq.Connect()
	if err != nil {
		log.Panic("Failed to connect to RabbitMQ ❌")
	}
	defer conn.Close()

	app := Config{
		rmqConn: conn,
	}

	log.Printf("Starting broker service on port: %s\n", webPort)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}
