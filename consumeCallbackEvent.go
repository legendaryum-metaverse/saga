package saga

import (
	"encoding/json"
	"fmt"
	"log"
	"slices"
	"time"

	"github.com/google/uuid"
	"github.com/legendaryum-metaverse/saga/event"

	amqp "github.com/rabbitmq/amqp091-go"
)

type EventHandler struct {
	Channel *EventsConsumeChannel  `json:"channel"`
	Payload map[string]interface{} `json:"payload"`
}

func ParsePayload[T any](handlerPayload map[string]interface{}, data *T) *T {
	jsonData, err := json.Marshal(handlerPayload)
	if err != nil {
		panic(fmt.Errorf("failed to marshal payload: %w", err))
	}
	if err = json.Unmarshal(jsonData, &data); err != nil {
		panic(fmt.Errorf("failed to unmarshal payload: %w", err))
	}
	if err = validate.Struct(data); err != nil {
		panic(fmt.Errorf("invalid payload: %w", err))
	}
	return data
}

// ParseEventPayload also works, but you need to pass a reference to the variable
// and is not type safe to assure that, as the type is: any
// Works:
// var eventPayload1 saga.SocialNewUserPayload   // or a pointer *saga.SocialNewUserPayload
// ------------------------->key, pass the reference<-----------------//
// handler.ParsePayload(&eventPayload1)
//
// It does not work:
// handler.ParsePayload(eventPayload1).
func (e *EventHandler) ParseEventPayload(data any) {
	body, err := json.Marshal(e.Payload)
	if err != nil {
		panic(err)
	}
	if err = json.Unmarshal(body, &data); err != nil {
		panic(err)
	}
}

// eventCallback handles the consumption and processing of microservice events.
func (t *Transactional) eventCallback(msg *amqp.Delivery, emitter *Emitter[EventHandler, event.MicroserviceEvent], queueName string) {
	if msg == nil {
		fmt.Println("Message not available")
		return
	}

	var eventPayload map[string]interface{}
	if err := json.Unmarshal(msg.Body, &eventPayload); err != nil {
		fmt.Printf("Error parsing message: %s\n", err)
		err = t.eventsChannel.Nack(msg.DeliveryTag, false, false)
		if err != nil {
			fmt.Println("Error negatively acknowledging message:", err)
			return
		}
		return
	}

	eventKey, err := findEventValues(msg.Headers)
	if err != nil {
		fmt.Println("Invalid header value: no valid event key found")
		err = t.eventsChannel.Nack(msg.DeliveryTag, false, false)
		if err != nil {
			fmt.Println("Error negatively acknowledging message:", err)
			return
		}
		return
	}
	if len(eventKey) != 1 {
		fmt.Println("More then one valid header, using the first one detected, that is because the payload is typed with a particular event")
	}

	eventType := string(eventKey[0])

	// Extract publisher microservice from message properties (set by the publisher)
	publisherMicroservice := msg.AppId
	if publisherMicroservice == "" {
		log.Printf("Warning: Message is missing AppId, using unknown as publisher_microservice")
		publisherMicroservice = "unknown"
	}

	// Extract or generate event_id from message properties
	eventID := msg.MessageId
	if eventID == "" {
		log.Printf("Warning: Message is missing MessageId, generating a new UUID v7 for event_id")
		eventID = uuid.Must(uuid.NewV7()).String()
	}
	// Emit audit.received event automatically when event is received (before processing)
	timestamp := uint64(time.Now().Unix())

	auditReceivedPayload := event.AuditReceivedPayload{
		PublisherMicroservice: publisherMicroservice,
		ReceiverMicroservice:  string(t.Microservice),
		ReceivedEvent:         eventType,
		ReceivedAt:            timestamp,
		QueueName:             queueName,
		EventID:               eventID,
	}
	go func() {
		// Emit the audit.received event (don't fail the main flow if audit fails)
		if auditErr := PublishAuditEvent(&auditReceivedPayload); auditErr != nil {
			log.Printf("Failed to emit audit.received event: %v", auditErr)
		}
	}()

	responseChannel := &EventsConsumeChannel{
		ConsumeChannel: &ConsumeChannel{
			channel:   t.eventsChannel,
			msg:       msg,
			queueName: queueName,
		},
		microservice:          string(t.Microservice),
		eventType:             eventType,
		publisherMicroservice: publisherMicroservice,
		eventID:               eventID,
	}

	emitter.Emit(eventKey[0], EventHandler{
		Payload: eventPayload,
		Channel: responseChannel,
	})
}

// findEventValues find all the MicroserviceEvent values in the headers.
func findEventValues(headers amqp.Table) ([]event.MicroserviceEvent, error) {
	var eventValues []event.MicroserviceEvent
	for _, value := range headers {
		if _, ok := value.(string); !ok {
			continue
		}
		val := event.MicroserviceEvent(value.(string))
		if slices.Contains(event.MicroserviceEventValues(), val) {
			eventValues = append(eventValues, val)
		}
	}
	if len(eventValues) == 0 {
		return nil, fmt.Errorf("no valid event key found")
	}
	return eventValues, nil
}
