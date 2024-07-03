package saga

import (
	"fmt"

	"github.com/go-playground/validator/v10"

	"go.uber.org/zap"

	"github.com/legendaryum-metaverse/saga/event"
	"github.com/legendaryum-metaverse/saga/micro"
	amqp "github.com/rabbitmq/amqp091-go"
)

// consume consumes messages from the queue and processes them.
func consume[T any, U comparable](e *Emitter[T, U], queueName string, logger *zap.Logger, cb func(*amqp.Delivery, *amqp.Channel, *Emitter[T, U], string, *zap.Logger)) error {
	channel, err := getConsumeChannel()
	if err != nil {
		return fmt.Errorf("failed to get channel: %w", err)
	}

	// TODO: is the same channel as the createConsumers channel, if for some reason the prefetch count is changed in one place, it will be changed in the other
	// TODO: solution, create a new channel
	err = channel.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		return fmt.Errorf("failed to set QoS: %w", err)
	}

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
		cb(&msg, channel, e, queueName, logger)
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
	RabbitUri    string
	Microservice micro.AvailableMicroservices
	Events       []event.MicroserviceEvent
	isReady      bool
	Logger       *zap.Logger
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
	Logger       *zap.Logger                  `validate:"required"`
}

func Config(opts *Opts) *Transactional {
	err := validate.Struct(opts)
	if err != nil {
		panic(fmt.Sprintf("Invalid options: %v", err))
	}

	if opts.Events == nil {
		opts.Events = []event.MicroserviceEvent{}
	}

	return &Transactional{
		RabbitUri:    opts.RabbitUri,
		Microservice: opts.Microservice,
		Events:       opts.Events,
		Logger:       opts.Logger,
	}
}

// ConnectToSagaCommandEmitter connects to the saga commands exchange and returns an emitter.
func (t *Transactional) ConnectToSagaCommandEmitter() *Emitter[CommandHandler, micro.StepCommand] {
	t.prepare()
	q := getQueueConsumer(t.Microservice)
	e := newEmitter[CommandHandler, micro.StepCommand]()

	err := createConsumers([]QueueConsumerProps{q})
	if err != nil {
		panic(err)
	}

	go func() {
		err = consume(e, q.QueueName, t.Logger, sagaCommandCallback)
		if err != nil {
			fmt.Println("Error consuming messages:", err)
		}
	}()

	return e
}

// ConnectToEvents connects to the events exchange and returns an emitter.
func (t *Transactional) ConnectToEvents() *Emitter[EventHandler, event.MicroserviceEvent] {
	t.prepare()
	queueName := fmt.Sprintf("%s_match_commands", t.Microservice)
	e := newEmitter[EventHandler, event.MicroserviceEvent]()

	err := createHeaderConsumer(queueName, t.Events)
	if err != nil {
		panic(err)
	}

	go func() {
		err = consume(e, queueName, t.Logger, eventCallback)
		if err != nil {
			fmt.Println("Error consuming messages:", err)
		}
	}()

	return e
}
