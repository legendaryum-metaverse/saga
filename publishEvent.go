package saga

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/legendaryum-metaverse/saga/event"
	amqp "github.com/rabbitmq/amqp091-go"
)

func (t *Transactional) PublishEvent(payload event.PayloadEvent) error {
	headerEvent := getEventObject(payload.Type())
	headersArgs := amqp.Table{
		"all-micro": "yes",
	}
	for k, v := range headerEvent {
		headersArgs[k] = v
	}

	err := t.sendChannel.ExchangeDeclare(string(MatchingExchange), "headers", true, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("failed to declare exchange %s: %w", string(MatchingExchange), err)
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = t.sendChannel.PublishWithContext(
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
