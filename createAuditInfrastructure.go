package saga

import (
	"fmt"
)

// createAuditLoggingResources creates audit logging infrastructure with direct exchange and separate queues.
// Uses direct exchange for efficient single-consumer delivery to audit microservice.
func (t *Transactional) createAuditLoggingResources() error {
	// Create direct exchange for audit events
	err := t.eventsChannel.ExchangeDeclare(
		string(AuditExchange), // name
		"direct",              // type
		true,                  // durable
		false,                 // auto-deleted
		false,                 // internal
		false,                 // no-wait
		nil,                   // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare audit exchange: %w", err)
	}

	// Create separate queue for audit.received events
	_, err = t.eventsChannel.QueueDeclare(
		string(AuditReceivedCommandsQ), // name
		true,                           // durable
		false,                          // delete when unused
		false,                          // exclusive
		false,                          // no-wait
		nil,                            // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare audit.received queue: %w", err)
	}

	// Create separate queue for audit.processed events
	_, err = t.eventsChannel.QueueDeclare(
		string(AuditProcessedCommandsQ), // name
		true,                            // durable
		false,                           // delete when unused
		false,                           // exclusive
		false,                           // no-wait
		nil,                             // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare audit.processed queue: %w", err)
	}

	// Create separate queue for audit.dead_letter events
	_, err = t.eventsChannel.QueueDeclare(
		string(AuditDeadLetterCommandsQ), // name
		true,                             // durable
		false,                            // delete when unused
		false,                            // exclusive
		false,                            // no-wait
		nil,                              // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare audit.dead_letter queue: %w", err)
	}

	// Bind each queue to its specific routing key
	err = t.eventsChannel.QueueBind(
		string(AuditReceivedCommandsQ), // queue name
		"audit.received",               // routing key
		string(AuditExchange),          // exchange
		false,                          // no-wait
		nil,                            // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to bind audit.received queue: %w", err)
	}

	err = t.eventsChannel.QueueBind(
		string(AuditProcessedCommandsQ), // queue name
		"audit.processed",               // routing key
		string(AuditExchange),           // exchange
		false,                           // no-wait
		nil,                             // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to bind audit.processed queue: %w", err)
	}

	err = t.eventsChannel.QueueBind(
		string(AuditDeadLetterCommandsQ), // queue name
		"audit.dead_letter",              // routing key
		string(AuditExchange),            // exchange
		false,                            // no-wait
		nil,                              // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to bind audit.dead_letter queue: %w", err)
	}

	return nil
}
