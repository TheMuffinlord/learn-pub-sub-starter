package main

import (
	"fmt"
	"learn-pub-sub-starter/internal/gamelogic"
	"learn-pub-sub-starter/internal/pubsub"
	"learn-pub-sub-starter/internal/routing"
	"log"

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

	gamelogic.PrintServerHelp()

	for {
		serverCmds := gamelogic.GetInput()
		if len(serverCmds) != 0 {
			switch serverCmds[0] {
			case "pause":
				fmt.Println("Pausing server.")
				err = pubsub.PublishJSON(rabbitChan, string(routing.ExchangePerilDirect), string(routing.PauseKey), routing.PlayingState{IsPaused: true})
				if err != nil {
					log.Fatalf("error publishing message to exchange: %v", err)

				}
				fmt.Println("Pause message published to exchange.")
			case "resume":
				fmt.Println("Resuming server.")
				err = pubsub.PublishJSON(rabbitChan, string(routing.ExchangePerilDirect), string(routing.PauseKey), routing.PlayingState{IsPaused: false})
				if err != nil {
					log.Fatalf("error publishing message to exchange: %v", err)

				}
				fmt.Println("Resume message published to exchange.")
			case "quit":
				fmt.Println("Exiting server.")
				return
			default:
				fmt.Println("I don't understand that command, sorry.")
			}
		}
	}
	/*
		// wait for ctrl+c
		signalChan := make(chan os.Signal, 1)
		signal.Notify(signalChan, os.Interrupt)
		<-signalChan
		fmt.Println("Closing server.")*/
}
