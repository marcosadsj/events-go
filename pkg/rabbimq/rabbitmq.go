package rabbimq

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func OpenChannel() (*amqp.Channel, error) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %s", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %s", err)
	}

	return ch, nil
}

func Consume(ch *amqp.Channel, out chan<- amqp.Delivery, queueName string) error {
	msgs, err := ch.Consume(
		queueName,     // queue
		"go-consumer", // consumer
		false,         // auto-ack
		false,         // exclusive
		false,         // no-local
		false,         // no-wait
		nil,           // args
	)
	if err != nil {
		return err
	}

	for d := range msgs {
		out <- d
	}

	return nil
}

func Publish(ch *amqp.Channel, exchange, routingKey string, body []byte) error {
	err := ch.Publish(
		exchange,   // exchange
		routingKey, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
		},
	)
	if err != nil {
		return err
	}

	return nil
}
