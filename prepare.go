package saga

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

// HealthCheck checks if the rabbitmq connection is alive and the queue exists.
// The queue to check is the microservice related to the saga commands or events.
func (t *Transactional) HealthCheck() error {
	if !t.isConnected {
		return fmt.Errorf("rabbitmq is not connected")
	}

	// pueden ser las dos o una de las dos
	if t.healthCheckQueue == "" {
		return fmt.Errorf("health check queue is not set")
	}
	channel, err := t.conn.Channel()

	defer func(channel *amqp.Channel) {
		if channel != nil {
			err := channel.Close()
			if err != nil {
				fmt.Println("channel close error")
			}
		}
	}(channel)
	if err != nil {
		fmt.Println("channel error")
		return err

	}

	_, err = channel.QueueDeclarePassive(t.healthCheckQueue, true, false, false, false, nil)
	if err != nil {
		fmt.Println("queue error")
		return err
	}

	return nil
}

// StopRabbitMQ closes the rabbitmq connection and channels.
func (t *Transactional) StopRabbitMQ() error {
	var err error
	if t.eventsChannel != nil {
		err = t.eventsChannel.Close()
	}
	if t.sagaChannel != nil {
		err = t.sagaChannel.Close()
	}
	if t.sendChannel != nil {
		err = t.sendChannel.Close()
	}
	if t.conn != nil {
		err = t.conn.Close()
	}
	return err
}

func (t *Transactional) notifyClose() {
	closeReceiver := make(chan *amqp.Error)
	t.conn.NotifyClose(closeReceiver)

	go func() {
		<-closeReceiver
		t.isConnected = false
	}()
}
