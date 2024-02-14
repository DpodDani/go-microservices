package event

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

const EXCHANGE_NAME string = "logs_topic"

func declareExchange(channel *amqp.Channel) error {
	return channel.ExchangeDeclare(
		EXCHANGE_NAME, // name
		"topic",       // type
		true,          // durable?
		false,         // auto-delete?
		false,         // internal
		false,         // no wait,
		nil,           // arguments
	)
}

func declareRandomQueue(channel *amqp.Channel) (amqp.Queue, error) {
	return channel.QueueDeclare(
		"",    // name
		false, // durable?
		false, // delete when unused?
		false, // exclusive?
		false, // noWait?
		nil,   // arguments
	)
}
