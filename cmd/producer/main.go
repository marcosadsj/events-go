package main

import "events-go/pkg/rabbimq"

func main() {

	ch, err := rabbimq.OpenChannel()

	if err != nil {
		panic(err)
	}
	defer ch.Close()

	rabbimq.Publish(ch, "amq.direct", "", []byte("Hello, World!"))

}
