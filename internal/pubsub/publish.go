package pubsub

import (
	"context"
	"encoding/json"

	amqp "github.com/rabbitmq/amqp091-go"
)

type SimpleQueueType struct {
	Durable   bool
	Transient bool
}

func PublishJSON[T any](ch *amqp.Channel, exchange, key string, val T) error {
	jsonData, err := json.Marshal(val)
	if err != nil {
		return err
	}

	if err := ch.PublishWithContext(
		context.Background(),
		exchange,
		key,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        jsonData,
		},
	); err != nil {
		return err
	}

	return nil
}

func DeclareAndBind(
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	queueType SimpleQueueType,
) (*amqp.Channel, amqp.Queue, error) {
	pubsubChan, err := conn.Channel()
	if err != nil {
		return nil, amqp.Queue{}, err
	}

	var queue amqp.Queue
	if queueType.Durable {
		queue, err = pubsubChan.QueueDeclare(
			queueName,
			true,
			false,
			false,
			false,
			amqp.Table{
				"x-dead-letter-exchange": "peril_dlx",
			},
		)
		if err != nil {
			return nil, amqp.Queue{}, err
		}
	}
	if queueType.Transient {
		queue, err = pubsubChan.QueueDeclare(
			queueName,
			false,
			true,
			true,
			false,
			amqp.Table{
				"x-dead-letter-exchange": "peril_dlx",
			},
		)
		if err != nil {
			return nil, amqp.Queue{}, err
		}
	}

	if err := pubsubChan.QueueBind(
		queueName,
		key,
		exchange,
		false,
		nil,
	); err != nil {
		return nil, amqp.Queue{}, err
	}

	return pubsubChan, queue, nil
}
