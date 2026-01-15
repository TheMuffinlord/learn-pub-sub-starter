package main

import (
	"fmt"
	"learn-pub-sub-starter/internal/gamelogic"
	"learn-pub-sub-starter/internal/pubsub"
	"learn-pub-sub-starter/internal/routing"
	"log"
	"strconv"

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

	publishCh, err := conn.Channel()
	if err != nil {
		log.Fatalf("Could not create channel: %v", err)
	}

	userName, err := gamelogic.ClientWelcome()
	if err != nil {
		log.Fatalf("Welcome Error: %v", err)
	}
	//_, queue, err := pubsub.DeclareAndBind(conn, routing.ExchangePerilDirect, "pause."+userName, "pause", "transient")
	/*if err != nil {
		log.Fatalf("Error binding channel and queue: %v", err)
	}
	fmt.Printf("Queue %v declared and bound successfully.\n", queue.Name)*/

	gs := gamelogic.NewGameState(userName)
	err = pubsub.SubscribeJSON(
		conn,
		routing.ExchangePerilDirect,
		routing.PauseKey+"."+gs.GetUsername(),
		routing.PauseKey, "transient", handlerPause(gs))
	if err != nil {
		log.Fatalf("Error subscribing to pause queue: %v", err)
	}
	err = pubsub.SubscribeJSON(conn,
		routing.ExchangePerilTopic,
		routing.ArmyMovesPrefix+"."+gs.GetUsername(),
		routing.ArmyMovesPrefix+".*",
		"transient", handleMove(gs, publishCh))
	if err != nil {
		log.Fatalf("Error subscribing to move queue: %v", err)
	}
	err = pubsub.SubscribeJSON(conn,
		routing.ExchangePerilTopic,
		routing.WarRecognitionsPrefix,
		routing.WarRecognitionsPrefix+".*",
		pubsub.QueueTypeDurable,
		handleWar(gs, publishCh),
	)
	if err != nil {
		log.Fatalf("Error subscribing to war channel: %v", err)
	}
	fmt.Println("Subscribed to all channels.")
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
				mv, err := gs.CommandMove(clientCmds)
				if err != nil {
					fmt.Printf("invalid move syntax: %s\n", err)
					continue
				}
				err = pubsub.PublishJSON(publishCh,
					routing.ExchangePerilTopic,
					routing.ArmyMovesPrefix+"."+mv.Player.Username,
					mv,
				)
				if err != nil {
					fmt.Printf("error: %v", err)
					continue
				}
				fmt.Printf("Moved %v unit(s) to %v\n", len(mv.Units), mv.ToLocation)
			case "status":
				gs.CommandStatus()
			case "help":
				gamelogic.PrintClientHelp()
			case "spam":
				if len(clientCmds) != 2 {
					fmt.Println("usage: spam <quantity>")
					continue
				}
				spamQuant, err := strconv.Atoi(clientCmds[1])
				if err != nil {
					fmt.Printf("error with spam: %v\n", err)
					fmt.Println("usage: spam <quantity>")
					continue
				}
				fmt.Println("spamming queue...")
				for i := 0; i < spamQuant; i++ {
					spamLog := gamelogic.GetMaliciousLog()
					err = pubsub.PublishGameLog(publishCh, userName, spamLog)
					if err != nil {
						fmt.Printf("error during spam: %v\n", err)
						break
					}
				}
				fmt.Printf("spam message published %v times.\n", spamQuant)

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
