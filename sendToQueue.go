package saga

import (
	"context"
	"encoding/json"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func (m *MicroserviceConsumeChannel) sendToQueue(queueName Queue, step SagaStep) error {
	err := send(m.channel, string(queueName), step)
	if err != nil {
		return err
	}
	return nil
}

func send(channel *amqp.Channel, queueName string, payload interface{}) error {
	_, err := channel.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		return err
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = channel.PublishWithContext(ctx, "", queueName, false, false, amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType:  "application/json",
		Body:         body,
	})
	if err != nil {
		return err
	}

	return nil
}
