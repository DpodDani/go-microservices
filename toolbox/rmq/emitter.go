package rmq

import (
	"context"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Emitter struct {
	conn *amqp.Connection
}

func NewEmitter(conn *amqp.Connection) (Emitter, error) {
	e := Emitter{
		conn: conn,
	}

	err := e.setup()
	if err != nil {
		return Emitter{}, err
	}

	return e, nil
}

func (e *Emitter) setup() error {
	ch, err := e.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()
	return declareExchange(ch)
}

func (e *Emitter) Push(event string, severity string) error {
	ch, err := e.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	log.Printf("Pushing message to channel: %s", event)
	err = ch.PublishWithContext(
		context.TODO(),
		EXCHANGE_NAME, // exchange
		severity,      // key
		false,         // mandatory?
		false,         // immediate
		amqp.Publishing{
			ContentType: "plain/text",
			Body:        []byte(event),
		},
	)
	if err != nil {
		return err
	}

	return nil
}
