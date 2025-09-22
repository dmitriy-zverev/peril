package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"

	"github.com/dmitriy-zverev/peril/internal/gamelogic"
	"github.com/dmitriy-zverev/peril/internal/pubsub"
	"github.com/dmitriy-zverev/peril/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	fmt.Println("Starting Peril client...")

	connString := "amqp://guest:guest@localhost:5672/"
	conn, err := amqp.Dial(connString)
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	fmt.Printf("Connection to RabbitMQ is successful at %s\n", connString)

	username, err := gamelogic.ClientWelcome()
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}

	_, _, err = pubsub.DeclareAndBind(
		conn,
		routing.ExchangePerilDirect,
		routing.PauseKey+"."+username,
		routing.PauseKey,
		pubsub.SimpleQueueType{Transient: true},
	)
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}

	gamestate := gamelogic.NewGameState(username)
	for {
		input := gamelogic.GetInput()
		if len(input) < 1 {
			continue
		}

		stop, err := parseClientCommand(input, gamestate)
		if err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}
		if stop {
			break
		}
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
	fmt.Println("\nShutting down...")
}

func parseClientCommand(
	cmds []string,
	gamestate *gamelogic.GameState,
) (bool, error) {
	switch strings.ToLower(cmds[0]) {
	case "spawn":
		if err := gamestate.CommandSpawn(cmds); err != nil {
			return false, err
		}
	case "move":
		move, err := gamestate.CommandMove(cmds)
		if err != nil {
			return false, err
		}

		movedUnitsRanks := []string{}
		for _, unit := range move.Units {
			movedUnitsRanks = append(movedUnitsRanks, string(unit.Rank))
		}

		fmt.Printf(
			"%s moved %v to %v\n",
			gamestate.Player.Username,
			strings.Join(movedUnitsRanks, ","),
			move.ToLocation,
		)
	case "status":
		gamestate.CommandStatus()
	case "spam":
		fmt.Println("Spamming not allowed yet!")
	case "quit":
		gamelogic.PrintQuit()
		return true, nil
	default:
		fmt.Printf("Don't know command '%s'\n", cmds[0])
	}

	return false, nil
}
