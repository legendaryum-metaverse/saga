package saga

import "github.com/legendaryum-metaverse/saga/micro"

type Status string

const (
	Success Status = "success"
	Failure Status = "failure"
	Sent    Status = "sent"
	Pending Status = "pending"
)

type SagaStep struct {
	Microservice    micro.AvailableMicroservices `json:"microservice"`
	Command         string                       `json:"command"`
	Status          Status                       `json:"status"`
	SagaID          int                          `json:"sagaId"`
	Payload         map[string]interface{}       `json:"payload"`
	PreviousPayload map[string]interface{}       `json:"previousPayload"`
	IsCurrentStep   bool                         `json:"isCurrentStep"`
}
