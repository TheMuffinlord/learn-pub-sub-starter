package main

import (
	"fmt"
	"learn-pub-sub-starter/internal/gamelogic"
	"learn-pub-sub-starter/internal/pubsub"
	"learn-pub-sub-starter/internal/routing"

	amqp "github.com/rabbitmq/amqp091-go"
)

//probably remember these, seems like you'll have to make more

func handlerPause(gs *gamelogic.GameState) func(routing.PlayingState) pubsub.AckType {
	return func(rps routing.PlayingState) pubsub.AckType {
		defer fmt.Print("> ")
		//fmt.Println("attempting to launch pause handler.")
		gs.HandlePause(rps)
		return pubsub.AckTypeAck
	}
}

func handleMove(gs *gamelogic.GameState, pubCh *amqp.Channel) func(gamelogic.ArmyMove) pubsub.AckType {
	return func(move gamelogic.ArmyMove) pubsub.AckType {
		defer fmt.Print("> ")
		outcome := gs.HandleMove(move)
		switch outcome {
		case gamelogic.MoveOutcomeSafe:
			return pubsub.AckTypeAck
		case gamelogic.MoveOutcomeMakeWar:
			err := pubsub.PublishJSON(pubCh, routing.ExchangePerilTopic, routing.WarRecognitionsPrefix+"."+gs.GetUsername(), gamelogic.RecognitionOfWar{
				Attacker: move.Player,
				Defender: gs.GetPlayerSnap(),
			})
			if err != nil {
				fmt.Printf("error: %s\n", err)
				return pubsub.AckTypeNackRequeue
			}
			return pubsub.AckTypeAck
		default:
			return pubsub.AckTypeNackDiscard
		}
	}
}

func handleWar(gs *gamelogic.GameState, pubCh *amqp.Channel) func(gamelogic.RecognitionOfWar) pubsub.AckType {
	return func(rw gamelogic.RecognitionOfWar) pubsub.AckType {
		defer fmt.Print("> ")
		outcome, winner, loser := gs.HandleWar(rw)
		switch outcome {
		case gamelogic.WarOutcomeNotInvolved:
			return pubsub.AckTypeNackRequeue
		case gamelogic.WarOutcomeNoUnits:
			return pubsub.AckTypeNackDiscard
		case gamelogic.WarOutcomeOpponentWon:
			logMsg := winner + " won a war against " + loser
			err := pubsub.PublishGameLog(gs, pubCh, logMsg)
			if err != nil {
				fmt.Printf("error logging war outcome: %v", err)
			}
			return pubsub.AckTypeAck
		case gamelogic.WarOutcomeYouWon:
			logMsg := winner + " won a war against " + loser
			err := pubsub.PublishGameLog(gs, pubCh, logMsg)
			if err != nil {
				fmt.Printf("error logging war outcome: %v", err)
			}
			return pubsub.AckTypeAck
		case gamelogic.WarOutcomeDraw:
			logMsg := "A war between " + winner + " and " + loser + " resulted in a draw"
			err := pubsub.PublishGameLog(gs, pubCh, logMsg)
			if err != nil {
				fmt.Printf("error logging war outcome: %v", err)
			}
			return pubsub.AckTypeAck
		default:
			fmt.Println("error: unsupported war outcome")
			return pubsub.AckTypeNackDiscard
		}
	}
}
