package pubsub

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

type SimpleQueueType string

const (
	QueueTypeDurable   SimpleQueueType = "durable"
	QueueTypeTransient SimpleQueueType = "transient"
)

func DeclareAndBind(
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	queueType SimpleQueueType, // an enum to represent "durable" or "transient"
) (*amqp.Channel, amqp.Queue, error) {
	dbChannel, err := conn.Channel()

	if err != nil {
		return dbChannel, amqp.Queue{}, err
	}
	var isDurable, isTransient bool
	switch queueType {
	case QueueTypeDurable:
		isDurable = true
		isTransient = false
	case QueueTypeTransient:
		isDurable = false
		isTransient = true
	}
	dbQueue, err := dbChannel.QueueDeclare(queueName, isDurable, isTransient, isTransient, false, nil)
	if err != nil {
		return dbChannel, dbQueue, err
	}
	err = dbChannel.QueueBind(queueName, key, exchange, false, nil)
	if err != nil {
		return dbChannel, dbQueue, err
	}
	return dbChannel, dbQueue, nil
}
