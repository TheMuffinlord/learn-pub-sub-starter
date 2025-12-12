package main

import (
	"fmt"
	"learn-pub-sub-starter/internal/gamelogic"
	"learn-pub-sub-starter/internal/pubsub"
	"log"
	"os"
	"os/signal"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	fmt.Println("Starting Peril client...")
	const connectURL = "amqp://guest:guest@localhost:5672/"
	conn, err := amqp.Dial(connectURL)
	if err != nil {
		log.Fatalf("error creating new connection: %v", err)
	}
	defer conn.Close()
	fmt.Println("Connection to  server successful.")

	userName, err := gamelogic.ClientWelcome()
	if err != nil {
		log.Fatalf("Welcome Error: %v", err)
	}
	_, queue, err := pubsub.DeclareAndBind(conn, "peril_direct", "pause."+userName, "pause", "transient")
	if err != nil {
		log.Fatalf("Error binding channel and queue: %v", err)
	}
	fmt.Printf("Queue %v declared and bound successfully.\n", queue.Name)

	// wait for ctrl+c
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
	fmt.Println("RabbitMQ connection closed.")
}
