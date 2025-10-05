package saga

import (
	"fmt"
)

type (
	StepHashId string
	Occurrence int
)

type MicroserviceConsumeChannel struct {
	*ConsumeChannel
	step SagaStep
}

type NextStepPayload = map[string]interface{}

func (m *MicroserviceConsumeChannel) AckMessage(payloadForNextStep NextStepPayload) {
	m.step.Status = Success
	previousPayload := m.step.PreviousPayload
	metaData := make(map[string]interface{})

	for key, value := range previousPayload {
		if len(key) > 2 && key[:2] == "__" {
			metaData[key] = value
		}
	}

	for key, value := range payloadForNextStep {
		metaData[key] = value
	}

	m.step.Payload = metaData
	// Para que este micro pueda realizar pasos del saga y realizar commence_saga ops las queue's deben existir, no es responsabilidad
	// de los micros crear estos recursos, el micro "transactional" debe crear estos recursos -> "queue.CommenceSaga" en commenceSagaListener
	// y "queue.ReplyToSaga" en startGlobalSagaStepListener
	err := m.sendToQueue(ReplyToSagaQ, m.step)
	if err != nil {
		// TODO: reenqueue message o manejar mejor el error
		return
	}

	err = m.channel.Ack(m.msg.DeliveryTag, false)
	if err != nil {
		// TODO: reenqueue message
		fmt.Println("Error acknowledging message:", err)
	}
}
