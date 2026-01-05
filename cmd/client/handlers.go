package main

import (
	"fmt"
	"learn-pub-sub-starter/internal/gamelogic"
	"learn-pub-sub-starter/internal/pubsub"
	"learn-pub-sub-starter/internal/routing"
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

func handleMove(gs *gamelogic.GameState) func(gamelogic.ArmyMove) pubsub.AckType {
	return func(move gamelogic.ArmyMove) pubsub.AckType {
		defer fmt.Print("> ")
		outcome := gs.HandleMove(move)
		switch outcome {
		case gamelogic.MoveOutcomeSafe, gamelogic.MoveOutcomeMakeWar:
			return pubsub.AckTypeAck
		default:
			return pubsub.AckTypeNackDiscard
		}
	}
}
