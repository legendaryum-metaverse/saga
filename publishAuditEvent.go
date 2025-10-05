package saga

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/legendaryum-metaverse/saga/event"
	amqp "github.com/rabbitmq/amqp091-go"
)

// publishAuditEvent publishes audit events to the direct audit exchange.
// Uses the event type as routing key for flexible audit event routing.
func publishAuditEvent(payload event.PayloadEvent) error {
	channel, err := getSendChannel()
	if err != nil {
		return fmt.Errorf("error getting send channel: %w", err)
	}

	// Use the event type as routing key for flexible audit event routing
	eventType := payload.Type()
	routingKey := string(eventType) // "audit.received", "audit.processed", or "audit.dead_letter"

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal audit payload: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = channel.PublishWithContext(
		ctx,
		string(AuditExchange), // exchange
		routingKey,            // routing key
		false,                 // mandatory
		false,                 // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent, // persistent
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish audit event: %w", err)
	}

	return nil
}

// PublishAuditReceivedEvent publishes audit.received events - convenience wrapper.
func PublishAuditReceivedEvent(payload *event.AuditReceivedPayload) error {
	return publishAuditEvent(payload)
}

// PublishAuditProcessedEvent publishes audit.processed events - convenience wrapper.
func PublishAuditProcessedEvent(payload *event.AuditProcessedPayload) error {
	return publishAuditEvent(payload)
}

// PublishAuditDeadLetterEvent publishes audit.dead_letter events - convenience wrapper.
func PublishAuditDeadLetterEvent(payload *event.AuditDeadLetterPayload) error {
	return publishAuditEvent(payload)
}
