package saga

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/legendaryum-metaverse/saga/micro"
	amqp "github.com/rabbitmq/amqp091-go"
)

type CommandHandler struct {
	Channel *MicroserviceConsumeChannel `json:"channel"`
	Payload map[string]interface{}      `json:"payload"`
	SagaID  int                         `json:"sagaId"`
	// OperationID is the SIPLEI operation the saga step belongs to. Empty when
	// the publisher did not set it, expected during the migration window.
	OperationID string `json:"operationId"`
	// Ctx carries OperationID for downstream publishes and gRPC calls.
	Ctx context.Context `json:"-"`
}

func (t *Transactional) sagaCommandCallback(msg *amqp.Delivery, e *Emitter[CommandHandler, micro.StepCommand], queueName string) {
	if msg == nil {
		fmt.Println("NO MSG AVAILABLE")
		return
	}

	var currentStep SagaStep
	err := json.Unmarshal(msg.Body, &currentStep)
	if err != nil {
		fmt.Println("ERROR PARSING MSG", err)
		err = t.sagaChannel.Nack(msg.DeliveryTag, false, false)
		if err != nil {
			fmt.Println("Error negatively acknowledging message:", err)
			return
		}
		return
	}

	operationID := operationFromHeaders(msg.Headers)

	responseChannel := &MicroserviceConsumeChannel{
		step:        currentStep,
		operationID: operationID,
		ConsumeChannel: &ConsumeChannel{
			channel:   t.sagaChannel,
			msg:       msg,
			queueName: queueName,
		},
	}

	if operationID == "" {
		reportMissingOperation(string(t.Microservice), string(currentStep.Command))
	}

	e.Emit(currentStep.Command, CommandHandler{
		Channel:     responseChannel,
		Payload:     currentStep.PreviousPayload,
		SagaID:      currentStep.SagaID,
		OperationID: operationID,
		Ctx:         WithOperation(context.Background(), operationID),
	})
}
