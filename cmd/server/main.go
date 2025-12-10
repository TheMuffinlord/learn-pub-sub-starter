package main

import (
	"fmt"
	"learn-pub-sub-starter/internal/pubsub"
	"learn-pub-sub-starter/internal/routing"
	"log"
	"os"
	"os/signal"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	fmt.Println("Starting Peril server...")
	const connectURL = "amqp://guest:guest@localhost:5672/"
	conn, err := amqp.Dial(connectURL)
	if err != nil {
		log.Fatalf("error creating new connection: %v", err)
	}
	defer conn.Close()
	fmt.Println("Connection to new server successful.")

	rabbitChan, err := conn.Channel()
	if err != nil {
		log.Fatalf("error opening connection channel: %v", err)
	}
	fmt.Println("Connection channel opened.")
	err = pubsub.PublishJSON(rabbitChan, string(routing.ExchangePerilDirect), string(routing.PauseKey), routing.PlayingState{IsPaused: true})
	if err != nil {
		log.Fatalf("erorr publishing message to exchange: %v", err)
	}
	fmt.Println("Message published to exchange.")
	// wait for ctrl+c
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
	fmt.Println("Closing server.")
}
