package pubsub

import (
	"encoding/json"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type AckType int

const (
	AckTypeAck AckType = iota
	AckTypeNackRequeue
	AckTypeNackDiscard
)

func SubscribeJSON[T any](
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	queueType SimpleQueueType, // an enum to represent "durable" or "transient"
	handler func(T) AckType,
) error {
	subChannel, subQueue, err := DeclareAndBind(conn, exchange, queueName, key, queueType)
	if err != nil {
		return fmt.Errorf("error declaring and binding queue: %v", err)
	}
	fmt.Println("Declared queue.")
	deliveries, err := subChannel.Consume(subQueue.Name, "", false, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("could not consume messages: %v", err)
	}

	fmt.Println("created delivery channel.")

	unmarshaller := func(data []byte) (T, error) {
		var target T
		err := json.Unmarshal(data, &target)
		return target, err
	}

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
			case AckTypeAck:
				msg.Ack(false)
				log.Println("message ack")
			case AckTypeNackDiscard:
				msg.Nack(false, false)
				log.Println("message nack(discard)")
			case AckTypeNackRequeue:
				msg.Nack(false, true)
				log.Println("message nack(requeue)")
			}
		}
	}()

	return nil
}
