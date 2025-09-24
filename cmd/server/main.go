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
	fmt.Println("Starting Peril server...")
	gamelogic.PrintServerHelp()

	connString := "amqp://guest:guest@localhost:5672/"
	conn, err := amqp.Dial(connString)
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	fmt.Printf("Connection to RabbitMQ is successful at %s\n", connString)

	pubChan, err := conn.Channel()
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}

	_, _, err = pubsub.DeclareAndBind(
		conn,
		routing.ExchangePerilTopic,
		routing.GameLogSlug,
		routing.GameLogSlug+"."+"*",
		pubsub.SimpleQueueDurable,
	)
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}

	err = pubsub.SubscribeGob(
		conn,
		routing.ExchangePerilTopic,
		routing.GameLogSlug,
		routing.GameLogSlug+".*",
		pubsub.SimpleQueueDurable,
		handlerLogs(),
	)
	if err != nil {
		log.Fatalf("could not starting consuming logs: %v", err)
	}

	for {
		input := gamelogic.GetInput()
		if len(input) < 1 {
			continue
		}

		if parseServerCommand(input[0], pubChan) {
			break
		}
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
	fmt.Println("\nShutting down...")
}

func parseServerCommand(cmd string, ch *amqp.Channel) bool {
	switch strings.ToLower(cmd) {
	case "pause":
		fmt.Println("Sending pause message")
		if err := pubsub.PublishJSON(
			ch,
			string(routing.ExchangePerilDirect),
			string(routing.PauseKey),
			routing.PlayingState{
				IsPaused: true,
			},
		); err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}
	case "resume":
		fmt.Println("Sending resume message")
		if err := pubsub.PublishJSON(
			ch,
			string(routing.ExchangePerilDirect),
			string(routing.PauseKey),
			routing.PlayingState{
				IsPaused: false,
			},
		); err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}
	case "quit":
		fmt.Println("Exiting...")
		return true
	default:
		fmt.Printf("Don't know command '%s'\n", cmd)
	}

	return false
}
