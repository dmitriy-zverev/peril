package pubsub

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Acktype int

type SimpleQueueType int

const (
	SimpleQueueDurable SimpleQueueType = iota
	SimpleQueueTransient
)

const (
	Ack Acktype = iota
	NackDiscard
	NackRequeue
)

func SubscribeJSON[T any](
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	queueType SimpleQueueType,
	handler func(T) Acktype,
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

			switch handler(dataStruct) {
			case Ack:
				delivery.Ack(false)
			case NackDiscard:
				delivery.Nack(false, false)
			case NackRequeue:
				delivery.Nack(false, true)
			}
		}
	}()

	return nil
}

func SubscribeGob[T any](
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	queueType SimpleQueueType,
	handler func(T) Acktype,
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
			if err := gob.NewDecoder(bytes.NewBuffer(delivery.Body)).Decode(&dataStruct); err != nil {
				fmt.Printf("%v\n", err)
				os.Exit(1)
			}

			switch handler(dataStruct) {
			case Ack:
				delivery.Ack(false)
			case NackDiscard:
				delivery.Nack(false, false)
			case NackRequeue:
				delivery.Nack(false, true)
			}
		}
	}()

	return nil
}
