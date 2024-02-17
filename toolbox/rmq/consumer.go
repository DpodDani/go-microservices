package rmq

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	conn *amqp.Connection
}

type Payload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func NewConsumer(conn *amqp.Connection) (Consumer, error) {
	consumer := Consumer{
		conn: conn,
	}

	err := consumer.setup()
	if err != nil {
		log.Println("Failed to setup consumer")
		return Consumer{}, err
	}

	return consumer, nil
}

func handlePayload(payload Payload) {
	switch payload.Name {
	case "log", "event":
		err := logEvent(payload)
		if err != nil {
			log.Println(err)
		}
	case "auth":
		// do some work here
	default:
		err := logEvent(payload)
		if err != nil {
			log.Println(err)
		}
	}
}

func logEvent(payload Payload) error {
	jsonData, _ := json.MarshalIndent(payload, "", "\t")

	request, err := http.NewRequest(http.MethodPost, "http://logger-service/log", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	if len(jsonData) > 0 {
		request.Header.Set("Content-Type", "application/json")
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusAccepted {
		return errors.New("log service failed to process log entry")
	}

	return nil
}

func (c *Consumer) setup() error {
	channel, err := c.conn.Channel()
	if err != nil {
		log.Println("Failed to get RMQ channel")
		return err
	}

	return declareExchange(channel)
}

func (c *Consumer) Listen(topics []string) error {
	ch, err := c.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	q, err := declareRandomQueue(ch)
	if err != nil {
		return err
	}

	for _, topic := range topics {
		err = ch.QueueBind(
			q.Name,
			topic,
			EXCHANGE_NAME,
			false,
			nil,
		)
		if err != nil {
			return err
		}
	}

	messages, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		return err
	}

	forever := make(chan bool)
	go func() {
		for msg := range messages {
			var payload Payload
			_ = json.Unmarshal(msg.Body, &payload)
			go handlePayload(payload)
		}
	}()

	log.Printf("Waiting for messages on [Exchange, Queue] --> [%s, %s]\n", EXCHANGE_NAME, q.Name)
	<-forever

	return nil
}
