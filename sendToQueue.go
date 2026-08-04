package saga

import (
	"context"
	"encoding/json"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func (m *MicroserviceConsumeChannel) sendToQueue(queueName Queue, step SagaStep) error {
	err := send(WithOperation(context.Background(), m.operationID), m.channel, string(queueName), step)
	if err != nil {
		return err
	}
	return nil
}

func send(ctx context.Context, channel *amqp.Channel, queueName string, payload interface{}) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	publishCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = channel.PublishWithContext(publishCtx, "", queueName, false, false, amqp.Publishing{
		Headers:      applyOperationHeader(ctx, nil),
		DeliveryMode: amqp.Persistent,
		ContentType:  "application/json",
		Body:         body,
	})
	if err != nil {
		return err
	}

	return nil
}
