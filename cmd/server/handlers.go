package main

import (
	"fmt"

	"github.com/dmitriy-zverev/peril/internal/gamelogic"
	"github.com/dmitriy-zverev/peril/internal/pubsub"
	"github.com/dmitriy-zverev/peril/internal/routing"
)

func handlerLogs() func(gamelog routing.GameLog) pubsub.Acktype {
	return func(gamelog routing.GameLog) pubsub.Acktype {
		defer fmt.Print("> ")

		err := gamelogic.WriteLog(gamelog)
		if err != nil {
			fmt.Printf("error writing log: %v\n", err)
			return pubsub.NackRequeue
		}
		return pubsub.Ack
	}
}
