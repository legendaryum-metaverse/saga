package saga

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/legendaryum-metaverse/saga/event"
	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	sendChannel       *amqp.Channel
	publishConnection *amqp.Connection
)

func getPublishConnection() (*amqp.Connection, error) {
	if publishConnection != nil {
		if publishConnection.IsClosed() {
			publishConnection = nil
			return getPublishConnection()
		}
		return publishConnection, nil
	}
	if RabbitUri == "" {
		return nil, fmt.Errorf("RabbitUri is not set")
	}
	var err error
	publishConnection, err = amqp.Dial(RabbitUri)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}
	return publishConnection, nil
}

func getSendChannel() (*amqp.Channel, error) {
	if sendChannel != nil {
		if sendChannel.IsClosed() {
			sendChannel = nil
			return getSendChannel()
		}
		return sendChannel, nil
	}
	conn, err := getPublishConnection()
	if err != nil {
		return nil, err
	}
	sendChannel, err = conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open a channel: %w", err)
	}
	return sendChannel, nil
}

func PublishEvent(payload event.PayloadEvent) error {
	channel, err := getSendChannel()
	if err != nil {
		return fmt.Errorf("error getting send channel: %w", err)
	}
	headerEvent := getEventObject(payload.Type())
	headersArgs := amqp.Table{
		"all-micro": "yes",
	}
	for k, v := range headerEvent {
		headersArgs[k] = v
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = channel.PublishWithContext(
		ctx,
		string(MatchingExchange),
		"",
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			Headers:      headersArgs,
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent,
		},
	)
	if err != nil {
		return fmt.Errorf("error publishing message: %w", err)
	}
	return nil
}
