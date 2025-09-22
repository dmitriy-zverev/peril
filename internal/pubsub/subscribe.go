package pubsub

import (
	"encoding/json"
	"fmt"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

func SubscribeJSON[T any](
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	queueType SimpleQueueType,
	handler func(T),
) error {
	ch, queue, err := DeclareAndBind(
		conn,
		exchange,
		queueName,
		key,
		queueType,
	)
	if err != nil {
		return err
	}

	delChan, err := ch.Consume(
		queue.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	go func() {
		for delivery := range delChan {
			var dataStruct T
			if err := json.Unmarshal(delivery.Body, &dataStruct); err != nil {
				fmt.Printf("%v\n", err)
				os.Exit(1)
			}

			handler(dataStruct)
			delivery.Ack(false)
		}
	}()

	return nil
}
