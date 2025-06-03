package main

import (
	"events-go/pkg/rabbimq"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {

	ch, err := rabbimq.OpenChannel()

	if err != nil {
		panic(err)
	}

	defer ch.Close()

	msgs := make(chan amqp.Delivery)

	go rabbimq.Consume(ch, msgs, "my-queue")

	for msg := range msgs {

		println("Received message:", string(msg.Body))

		msg.Ack(false)
	}

}
