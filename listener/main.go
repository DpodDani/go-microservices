package main

import (
	"log"
	"math"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	// try to connect to RabbitMQ
	conn, err := connect()
	if err != nil {
		log.Panic("Failed to connect to RabbitMQ ❌")
	}
	defer conn.Close()
	// start listening to messages
	// create consumer
	// watch queue and consume events (from a topic)
}

func connect() (*amqp.Connection, error) {
	var connection *amqp.Connection
	var err error

	baseDelay := 1 * time.Second
	retries := 0
	maxRetries := 5
	rmqDsn := "amqp://guest:guest@localhost"

	for {
		conn, err := amqp.Dial(rmqDsn)
		if err != nil {
			log.Printf("Failed to connecto RabbitMQ server at: %s\n", rmqDsn)
			retries++
		} else {
			log.Println("Connected to RMQ ✅")
			connection = conn
			break
		}

		if retries > maxRetries {
			log.Println("Reached max retry limit. Exiting...")
			break
		}

		backOffTime := time.Duration(math.Pow(2, float64(retries))) * baseDelay
		log.Printf("Backing off for %v...\n", backOffTime)
		time.Sleep(backOffTime)
	}

	return connection, err
}
