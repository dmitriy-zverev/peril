package main

import (
	"fmt"
	"log"
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
		log.Fatalf("could not dial to rabbitmq: %v", err)
	}
	defer conn.Close()

	fmt.Printf("Connection to RabbitMQ is successful at %s\n", connString)

	username, err := gamelogic.ClientWelcome()
	if err != nil {
		log.Fatalf("could not create user: %v", err)
	}

	publishCh, err := conn.Channel()
	if err != nil {
		log.Fatalf("could not create channel: %v", err)
	}

	gamestate := gamelogic.NewGameState(username)

	if err := pubsub.SubscribeJSON(
		conn,
		string(routing.ExchangePerilDirect),
		string(routing.PauseKey)+"."+username,
		string(routing.PauseKey),
		pubsub.SimpleQueueType{Transient: true},
		handlerPause(gamestate),
	); err != nil {
		log.Fatalf("could not subscribe: %v", err)
	}

	if err := pubsub.SubscribeJSON(
		conn,
		string(routing.ExchangePerilTopic),
		string(routing.ArmyMovesPrefix)+"."+username,
		string(routing.ArmyMovesPrefix)+"."+"*",
		pubsub.SimpleQueueType{Transient: true},
		handlerMove(gamestate, publishCh),
	); err != nil {
		log.Fatalf("could not subscribe: %v", err)
	}

	if err := pubsub.SubscribeJSON(
		conn,
		string(routing.ExchangePerilTopic),
		string(routing.WarRecognitionsPrefix),
		string(routing.WarRecognitionsPrefix)+".#",
		pubsub.SimpleQueueType{Durable: true},
		handlerWar(gamestate),
	); err != nil {
		log.Fatalf("could not subscribe: %v", err)
	}

	for {
		input := gamelogic.GetInput()
		if len(input) < 1 {
			continue
		}

		stop, err := parseClientCommand(input, gamestate, publishCh)
		if err != nil {
			fmt.Printf("%v\n", err)

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
	gs *gamelogic.GameState,
	ch *amqp.Channel,
) (bool, error) {
	switch strings.ToLower(cmds[0]) {
	case "spawn":
		if err := gs.CommandSpawn(cmds); err != nil {
			return false, err
		}
	case "move":
		move, err := gs.CommandMove(cmds)
		if err != nil {
			return false, err
		}

		if err := pubsub.PublishJSON(
			ch,
			string(routing.ExchangePerilTopic),
			string(routing.ArmyMovesPrefix)+"."+gs.Player.Username,
			move,
		); err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}
		fmt.Println("The move was published successfuly")

		movedUnitsRanks := []string{}
		for _, unit := range move.Units {
			movedUnitsRanks = append(movedUnitsRanks, string(unit.Rank))
		}

		fmt.Printf(
			"%s moved %v to %v\n",
			gs.Player.Username,
			strings.Join(movedUnitsRanks, ","),
			move.ToLocation,
		)
	case "status":
		gs.CommandStatus()
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
