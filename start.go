package saga

import (
	"fmt"

	"github.com/go-playground/validator/v10"

	"github.com/legendaryum-metaverse/saga/event"
	"github.com/legendaryum-metaverse/saga/micro"
	amqp "github.com/rabbitmq/amqp091-go"
)

// consume consumes messages from the queue and processes them.
func consume[T any, U comparable](e *Emitter[T, U], queueName string, channel *amqp.Channel, cb func(*amqp.Delivery, *amqp.Channel, *Emitter[T, U], string)) error {
	channelQ, err := channel.Consume(
		queueName,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to consume: %w", err)
	}

	for msg := range channelQ {
		cb(&msg, channel, e, queueName)
	}

	return nil
}

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
	sendChannel   *amqp.Channel
	isConnected   bool
	// healthCheckQueue is the queue to check if the rabbitmq connection is alive.
	// This queue is set in the creation of resources, consumers, queues, and exchanges.
	healthCheckQueue string
}

var validate *validator.Validate

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

	sendChannel, err := conn.Channel()
	if err != nil {
		panic(fmt.Sprintf("Failed to create sendChannel: %v", err))
	}

	t := &Transactional{
		Microservice: opts.Microservice,
		Events:       opts.Events,
		conn:         conn,
		sendChannel:  sendChannel,
	}
	t.notifyClose()
	t.isConnected = true
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
		err = consume(e, q.QueueName, t.sagaChannel, sagaCommandCallback)
		if err != nil {
			fmt.Println("Error consuming messages:", err)
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

	go func() {
		err = consume(e, queueName, t.eventsChannel, eventCallback)
		if err != nil {
			fmt.Println("Error consuming messages:", err)
		}
	}()

	return e
}
