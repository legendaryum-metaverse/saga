package saga

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
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

	// Generate UUID v7 for event tracking
	eventID := uuid.Must(uuid.NewV7()).String()

	// Get publisher microservice name from stored config
	config := GetStoredConfig()
	if config == nil {
		return fmt.Errorf("config not initialized - cannot determine publisher microservice")
	}
	publisherMicroservice := string(config.Microservice)

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
			MessageId:    eventID,               // UUID v7 for event tracking
			AppId:        publisherMicroservice, // Publisher microservice name
		},
	)
	if err != nil {
		return fmt.Errorf("error publishing message: %w", err)
	}
	timestamp := uint64(time.Now().Unix())
	auditPayload := event.AuditPublishedPayload{
		PublisherMicroservice: publisherMicroservice,
		PublishedEvent:        string(payload.Type()),
		PublishedAt:           timestamp,
		EventID:               eventID,
	}
	// Emit audit.published event (fire-and-forget - never fail the main flow)
	go func(auditPayload event.AuditPublishedPayload) {
		if auditErr := PublishAuditEvent(&auditPayload); auditErr != nil {
			log.Printf("Failed to emit audit.published event: %v", auditErr)
		}
	}(auditPayload)

	return nil
}
