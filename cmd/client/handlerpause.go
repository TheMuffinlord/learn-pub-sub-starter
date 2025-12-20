package main

import (
	"fmt"
	"learn-pub-sub-starter/internal/gamelogic"
	"learn-pub-sub-starter/internal/routing"
)

func handlerPause(gs *gamelogic.GameState) func(routing.PlayingState) {
	return func(rps routing.PlayingState) {
		defer fmt.Print("> ")
		//fmt.Println("attempting to launch pause handler.")
		gs.HandlePause(rps)
	}
}
