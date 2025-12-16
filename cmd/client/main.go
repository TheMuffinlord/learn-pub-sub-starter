package main

import (
	"fmt"
	"learn-pub-sub-starter/internal/gamelogic"
	"learn-pub-sub-starter/internal/pubsub"
	"log"

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

	gs := gamelogic.NewGameState(userName)

	for {
		clientCmds := gamelogic.GetInput()
		if len(clientCmds) != 0 {
			switch clientCmds[0] {
			case "spawn":
				err = gs.CommandSpawn(clientCmds)
				if err != nil {
					fmt.Printf("invalid spawn syntax: %v\n", err)
				}
			case "move":
				_, err := gs.CommandMove(clientCmds)
				if err != nil {
					fmt.Printf("invalid move syntax: %s\n", err)
				}
			case "status":
				gs.CommandStatus()
			case "help":
				gamelogic.PrintClientHelp()
			case "spam":
				fmt.Println("Spamming is not allowed yet!")
			case "quit":
				fmt.Println("exiting client.")
				return
			default:
				fmt.Println("I don't understand the command.")
			}
		}
	}

	/* wait for ctrl+c
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
	fmt.Println("RabbitMQ connection closed.")*/
}
