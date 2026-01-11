package pubsub

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type AckType int

const (
	Ack AckType = iota
	NackDiscard
	NackRequeue
)

func SubscribeJSON[T any](
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	queueType SimpleQueueType, // an enum to represent "durable" or "transient"
	handler func(T) AckType,
) error {
	return subscribe(
		conn,
		exchange,
		queueName,
		key,
		queueType,
		handler,
		func(data []byte) (T, error) {
			var target T
			err := json.Unmarshal(data, &target)
			return target, err
		},
	)
}

func SubscribeGob[T any](
	conn *amqp.Connection,
	exchange, queueName, key string,
	queueType SimpleQueueType,
	handler func(T) AckType,
) error {
	return subscribe(
		conn,
		exchange,
		queueName,
		key,
		queueType,
		handler,
		func(data []byte) (T, error) {
			buffer := bytes.NewBuffer(data)
			dec := gob.NewDecoder(buffer)
			var target T
			err := dec.Decode(&target)
			return target, err
		},
	)
}

func subscribe[T any](
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	simpleQueueType SimpleQueueType,
	handler func(T) AckType,
	unmarshaller func([]byte) (T, error),
) error {
	subChannel, subQueue, err := DeclareAndBind(conn, exchange, queueName, key, simpleQueueType)
	if err != nil {
		return fmt.Errorf("error declaring and binding queue: %v", err)
	}
	fmt.Println("Declared queue.")
	deliveries, err := subChannel.Consume(subQueue.Name, "", false, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("could not consume messages: %v", err)
	}

	fmt.Println("Created delivery channel.")

	go func() {
		defer subChannel.Close()
		for msg := range deliveries {
			target, err := unmarshaller(msg.Body)
			if err != nil {
				fmt.Printf("Could not unmarshal message: %v", err)
				continue
			}
			acktype := handler(target)
			switch acktype {
			case Ack:
				msg.Ack(false)
				log.Println("message ack")
			case NackDiscard:
				msg.Nack(false, false)
				log.Println("message nack(discard)")
			case NackRequeue:
				msg.Nack(false, true)
				log.Println("message nack(requeue)")
			}
		}
	}()

	return nil
}
