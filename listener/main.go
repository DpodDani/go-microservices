package main

import (
	"log"

	rmq "github.com/DpodDani/go-microservices-toolbox/rmq"
)

func main() {
	// try to connect to RabbitMQ
	conn, err := rmq.Connect()
	if err != nil {
		log.Panic("Failed to connect to RabbitMQ ❌")
	}
	defer conn.Close()
	// start listening to messages
	log.Println("Listening for and consuming RabbitMQ messages")
	// create consumer
	consumer, err := rmq.NewConsumer(conn)
	if err != nil {
		log.Panic("Failed to create a consumer")
	}
	// watch queue and consume events (from a topic)
	var topics []string = []string{"log.INFO", "log.WARNING", "log.ERROR"}
	err = consumer.Listen(topics)
	if err != nil {
		log.Panicf("Failed to listen to topics: %v\n", topics)
	}
}
