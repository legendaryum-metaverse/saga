package test

import (
	"testing"

	"github.com/legendaryum-metaverse/saga"
	"github.com/legendaryum-metaverse/saga/event"
	"github.com/legendaryum-metaverse/saga/micro"
	"github.com/stretchr/testify/suite"
)

type EventsTestSuite struct {
	suite.Suite
	e *saga.Emitter[saga.EventHandler, event.MicroserviceEvent]
}

func (suite *EventsTestSuite) SetupTest() {
	transactional := saga.Transactional{
		RabbitUri:    "amqp://rabbit:1234@localhost:5672",
		Microservice: micro.RoomInventory,
		Events: []event.MicroserviceEvent{
			event.SocialNewUserEvent,
			event.SocialBlockChatEvent,
		},
	}
	eventEmitter := transactional.ConnectToEvents()
	suite.e = eventEmitter
}

func (suite *EventsTestSuite) TestSubscribedEvents() {
	eventSocialNewUserReceived := make(chan *event.SocialNewUserPayload)
	eventSocialBlockChatReceived := make(chan *event.SocialBlockChatPayload)

	suite.e.On(event.SocialNewUserEvent, func(handler saga.EventHandler) {
		eventPayload := saga.ParseEventPayload(handler.Payload, &event.SocialNewUserPayload{})
		defer func(eventPayload *event.SocialNewUserPayload) {
			eventSocialNewUserReceived <- eventPayload
		}(eventPayload)
		handler.Channel.AckMessage()
	})
	suite.e.On(event.SocialBlockChatEvent, func(handler saga.EventHandler) {
		eventPayload := saga.ParseEventPayload(handler.Payload, &event.SocialBlockChatPayload{})
		defer func(eventPayload *event.SocialBlockChatPayload) {
			eventSocialBlockChatReceived <- eventPayload
		}(eventPayload)
		handler.Channel.AckMessage()
	})

	err := saga.PublishEvent(&event.SocialNewUserPayload{
		UserID: "1234",
	})
	suite.Require().NoError(err)
	err = saga.PublishEvent(&event.SocialBlockChatPayload{
		UserID:        "1234",
		UserToBlockID: "4321",
	})
	suite.Require().NoError(err)

	p1 := <-eventSocialNewUserReceived
	suite.Equal("1234", p1.UserID)

	p2 := <-eventSocialBlockChatReceived
	suite.Equal("1234", p2.UserID)
	suite.Equal("4321", p2.UserToBlockID)
}

func TestEventsTestSuite(t *testing.T) {
	suite.Run(t, new(EventsTestSuite))
}
