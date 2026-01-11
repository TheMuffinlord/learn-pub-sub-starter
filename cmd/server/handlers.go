package main

import (
	"fmt"
	"learn-pub-sub-starter/internal/gamelogic"
	"learn-pub-sub-starter/internal/pubsub"
	"learn-pub-sub-starter/internal/routing"
)

func handlerLogs() func(gameLog routing.GameLog) pubsub.AckType {
	return func(gameLog routing.GameLog) pubsub.AckType {
		defer fmt.Print("\n> ")

		err := gamelogic.WriteLog(gameLog)
		if err != nil {
			fmt.Printf("error writing log: %v/n", err)
			return pubsub.NackRequeue
		}
		return pubsub.Ack
	}
}
