package saga

import (
	"fmt"

	"github.com/go-playground/validator/v10"

	"github.com/legendaryum-metaverse/saga/event"
	"github.com/legendaryum-metaverse/saga/micro"
	amqp "github.com/rabbitmq/amqp091-go"
)

func getQueueName(microservice micro.AvailableMicroservices) string {
	return fmt.Sprintf("%s_saga_commands", microservice)
}

func getQueueConsumer(microservice micro.AvailableMicroservices) QueueConsumerProps {
	return QueueConsumerProps{
		QueueName: getQueueName(microservice),
		Exchange:  CommandsExchange,
	}
}

type Transactional struct {
	Microservice micro.AvailableMicroservices
	Events       []event.MicroserviceEvent

	conn          *amqp.Connection
	eventsChannel *amqp.Channel
	sagaChannel   *amqp.Channel
	isConnected   bool
	// healthCheckQueue is the queue to check if the rabbitmq connection is alive.
	// This queue is set in the creation of resources, consumers, queues, and exchanges.
	healthCheckQueue string
}

var (
	validate     *validator.Validate
	storedConfig *Transactional
)

// GetStoredConfig returns the stored Transactional configuration.
// This is used by publishEvent to access the current microservice name.
func GetStoredConfig() *Transactional {
	return storedConfig
}

func init() {
	validate = validator.New()
	err := validate.RegisterValidation("microservice", func(fl validator.FieldLevel) bool {
		ms := fl.Field().String()
		return micro.AvailableMicroservices(ms).IsValid()
	})
	if err != nil {
		panic(err)
	}
}

type Opts struct {
	RabbitUri    string                       `validate:"required,url"`
	Microservice micro.AvailableMicroservices `validate:"required,microservice"`
	Events       []event.MicroserviceEvent    `validate:"-"`
}

// RabbitUri is used for send channel connection.
var RabbitUri string

func Config(opts *Opts) *Transactional {
	err := validate.Struct(opts)
	if err != nil {
		panic(fmt.Sprintf("Invalid options: %v", err))
	}

	if opts.Events == nil {
		opts.Events = []event.MicroserviceEvent{}
	}

	conn, err := amqp.Dial(opts.RabbitUri)
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to RabbitMQ: %v", err))
	}

	t := &Transactional{
		Microservice: opts.Microservice,
		Events:       opts.Events,
		conn:         conn,
	}
	t.notifyClose()
	t.isConnected = true
	RabbitUri = opts.RabbitUri

	// Store configuration for access by publishEvent
	storedConfig = t

	return t
}

// ConnectToSagaCommandEmitter connects to the saga commands exchange and returns an emitter.
func (t *Transactional) ConnectToSagaCommandEmitter() *Emitter[CommandHandler, micro.StepCommand] {
	sagaChannel, err := t.conn.Channel()
	if err != nil {
		panic(fmt.Sprintf("Failed to create sagaChannel: %v", err))
	}
	t.sagaChannel = sagaChannel
	err = t.sagaChannel.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		panic(fmt.Sprintf("Failed to set QoS in sagaChannel: %v", err))
	}

	q := getQueueConsumer(t.Microservice)
	e := newEmitter[CommandHandler, micro.StepCommand]()

	err = t.createConsumers([]QueueConsumerProps{q})
	if err != nil {
		panic(err)
	}

	go func() {
		channelQ, err := t.sagaChannel.Consume(
			q.QueueName,
			"",
			false,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			fmt.Println("Error consuming messages:", err)
		}

		for msg := range channelQ {
			t.sagaCommandCallback(&msg, e, q.QueueName)
		}
	}()

	return e
}

// ConnectToEvents connects to the events exchange and returns an emitter.
func (t *Transactional) ConnectToEvents() *Emitter[EventHandler, event.MicroserviceEvent] {
	eventsChannel, err := t.conn.Channel()
	if err != nil {
		panic(fmt.Sprintf("Failed to create sagaChannel: %v", err))
	}
	t.eventsChannel = eventsChannel
	err = t.eventsChannel.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		panic(fmt.Sprintf("Failed to set QoS in eventsChannel: %v", err))
	}

	queueName := fmt.Sprintf("%s_match_commands", t.Microservice)
	e := newEmitter[EventHandler, event.MicroserviceEvent]()

	err = t.createHeaderConsumer(queueName, t.Events)
	if err != nil {
		panic(err)
	}

	// Create audit logging resources - this feature is related only to "events"
	err = t.createAuditLoggingResources()
	if err != nil {
		panic(fmt.Sprintf("Failed to create audit logging resources: %v", err))
	}

	go func() {
		channelQ, err := t.eventsChannel.Consume(
			queueName,
			"",
			false,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			fmt.Println("Error consuming messages:", err)
		}

		for msg := range channelQ {
			t.eventCallback(&msg, e, queueName)
		}
	}()

	return e
}

// HealthCheck checks if the rabbitmq connection is alive and the queue exists.
// The queue to check is the microservice related to the saga commands or events.
func (t *Transactional) HealthCheck() error {
	if !t.isConnected {
		return fmt.Errorf("rabbitmq is not connected")
	}

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
