package pubsub

import (
	"bytes"
	"context"
	"encoding/gob"
	"learn-pub-sub-starter/internal/routing"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func PublishGob[T any](ch *amqp.Channel, exchange, key string, val T) error {
	//pubByte, err := json.Marshal(val)
	var logfile bytes.Buffer
	enc := gob.NewEncoder(&logfile)
	err := enc.Encode(val)
	if err != nil {
		return err
	}
	err = ch.PublishWithContext(context.Background(), exchange, key, false, false, amqp.Publishing{ContentType: "application/gob", Body: logfile.Bytes()})
	if err != nil {
		return err
	}
	return nil
}

func PublishGameLog(ch *amqp.Channel, username, logText string) error {
	pubLog := routing.GameLog{
		CurrentTime: time.Now(),
		Message:     logText,
		Username:    username,
	}
	err := PublishGob(ch, routing.ExchangePerilTopic, routing.GameLogSlug+"."+username, pubLog)
	return err
}
