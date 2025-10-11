package saga

import (
	"fmt"
	"log"
	"time"

	"github.com/legendaryum-metaverse/saga/event"
)

type EventsConsumeChannel struct {
	*ConsumeChannel
	microservice          string
	eventType             string
	publisherMicroservice string
	eventID               string
}

func (m *EventsConsumeChannel) AckMessage() {
	err := m.channel.Ack(m.msg.DeliveryTag, false)
	if err != nil {
		fmt.Println("error acknowledging message: %w", err)
		return
	}
	// Emit audit.processed event automatically
	timestamp := uint64(time.Now().Unix())

	auditPayload := &event.AuditProcessedPayload{
		PublisherMicroservice: m.publisherMicroservice,
		ProcessorMicroservice: m.microservice,
		ProcessedEvent:        m.eventType,
		ProcessedAt:           timestamp,
		QueueName:             m.queueName,
		EventID:               m.eventID,
	}
	go func(auditPayload *event.AuditProcessedPayload) {
		// Emit the audit event using the direct exchange method
		if err = PublishAuditEvent(auditPayload); err != nil {
			// Log the error but don't fail the ack operation
			log.Printf("Failed to emit audit.processed event: %v", err)
		}
	}(auditPayload)
}

// NackWithDelay wraps the base method and emits audit.dead_letter events.
func (m *EventsConsumeChannel) NackWithDelay(delay time.Duration, maxRetries int32) (int32, time.Duration, error) {
	count, duration, err := m.ConsumeChannel.NackWithDelay(delay, maxRetries)
	rc := uint32(count)
	// Emit audit.dead_letter event automatically
	timestamp := uint64(time.Now().Unix())

	auditPayload := &event.AuditDeadLetterPayload{
		PublisherMicroservice: m.publisherMicroservice,
		RejectorMicroservice:  m.microservice,
		RejectedEvent:         m.eventType,
		RejectedAt:            timestamp,
		QueueName:             m.queueName,
		RejectionReason:       "delay",
		RetryCount:            &rc,
		EventID:               m.eventID,
	}
	go func(auditPayload *event.AuditDeadLetterPayload) {
		// Emit the audit event (don't fail if audit fails)
		if auditErr := PublishAuditEvent(auditPayload); auditErr != nil {
			log.Printf("Failed to emit audit.dead_letter event: %v", auditErr)
		}
	}(auditPayload)

	return count, duration, err
}

// NackWithFibonacciStrategy wraps the base method and emits audit.dead_letter events.
func (m *EventsConsumeChannel) NackWithFibonacciStrategy(maxOccurrence, maxRetries int32) (int32, time.Duration, int32, error) {
	count, duration, occurrence, err := m.ConsumeChannel.NackWithFibonacciStrategy(maxOccurrence, maxRetries)
	rc := uint32(count)
	// Emit audit.dead_letter event automatically
	timestamp := uint64(time.Now().Unix())

	auditPayload := &event.AuditDeadLetterPayload{
		PublisherMicroservice: m.publisherMicroservice,
		RejectorMicroservice:  m.microservice,
		RejectedEvent:         m.eventType,
		RejectedAt:            timestamp,
		QueueName:             m.queueName,
		RejectionReason:       "fibonacci_strategy",
		RetryCount:            &rc,
		EventID:               m.eventID,
	}
	go func(auditPayload *event.AuditDeadLetterPayload) {

		// Emit the audit event (don't fail if audit fails)
		if auditErr := PublishAuditEvent(auditPayload); auditErr != nil {
			log.Printf("Failed to emit audit.dead_letter event: %v", auditErr)
		}
	}(auditPayload)

	return count, duration, occurrence, err
}
