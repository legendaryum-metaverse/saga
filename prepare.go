package saga

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	rabbitMQURL    string
	rabbitMQConn   *amqp.Connection
	consumeChannel *amqp.Channel
	isConnected    bool
)

func saveURI(uri string) {
	if rabbitMQURL == "" {
		rabbitMQURL = uri
	}
}

func getRabbitMQConn() (*amqp.Connection, error) {
	if rabbitMQConn == nil {
		conn, err := amqp.Dial(rabbitMQURL)
		if err != nil {
			return nil, err
		}
		rabbitMQConn = conn
		isConnected = true
	}
	return rabbitMQConn, nil
}

func getConsumeChannel() (*amqp.Channel, error) {
	if consumeChannel == nil {
		conn, err := getRabbitMQConn()
		if err != nil {
			return nil, err
		}
		channel, err := conn.Channel()
		if err != nil {
			return nil, err
		}
		consumeChannel = channel
	}
	return consumeChannel, nil
}

// Useful to prepare the connection in publish events
// func Prepare(url string) {
//	saveURI(url)
//	_, err := getRabbitMQConn()
//	if err != nil {
//		panic(err)
//	}
//	_, err = getConsumeChannel()
//	if err != nil {
//		panic(err)
//	}
//	notifyClose()
//}

func (t *Transactional) prepare() {
	if t.isReady {
		return
	}
	saveURI(t.RabbitUri)
	_, err := getRabbitMQConn()
	if err != nil {
		panic(err)
	}
	_, err = getConsumeChannel()
	if err != nil {
		panic(err)
	}
	t.isReady = true
	notifyClose()
}

// healthCheckQueue is the queue to check if the rabbitmq connection is alive.
// This queue is set in the creation of resources, consumers, queues, and exchanges.
var healthCheckQueue string

// HealthCheck checks if the rabbitmq connection is alive and the queue exists.
// The queue to check is the microservice related to the saga commands or events.
func HealthCheck() error {
	if !isConnected {
		return fmt.Errorf("rabbitmq is not connected")
	}
	if healthCheckQueue == "" {
		return fmt.Errorf("health check queue is not set")
	}
	channel, err := rabbitMQConn.Channel()
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
	_, err = channel.QueueDeclarePassive(healthCheckQueue, true, false, false, false, nil)
	if err != nil {
		fmt.Println("queue error")
		return err
	}
	return nil
}

func closeConsumeChannel() error {
	if consumeChannel != nil {
		err := consumeChannel.Close()
		consumeChannel = nil
		return err
	}
	return nil
}

func closeRabbitMQConn() error {
	if rabbitMQConn != nil {
		err := rabbitMQConn.Close()
		rabbitMQConn = nil
		rabbitMQURL = ""
		return err
	}
	return nil
}

func StopRabbitMQ() error {
	err := closeConsumeChannel()
	if err != nil {
		return err
	}
	err = closeSendChannel()
	if err != nil {
		return err
	}
	err = closeRabbitMQConn()
	if err != nil {
		return err
	}
	return nil
}

func notifyClose() {
	closeReceiver := make(chan *amqp.Error)
	rabbitMQConn.NotifyClose(closeReceiver)

	go func() {
		<-closeReceiver
		isConnected = false
	}()
}
