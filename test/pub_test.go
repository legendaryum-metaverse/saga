package test

import (
	"testing"
	"time"

	"github.com/legendaryum-metaverse/saga"
	"github.com/legendaryum-metaverse/saga/event"
	"github.com/legendaryum-metaverse/saga/micro"
	"github.com/stretchr/testify/suite"
)

type EventsTestSuite struct {
	suite.Suite
	e *saga.Emitter[saga.EventHandler, event.MicroserviceEvent]
	t *saga.Transactional
}

func (suite *EventsTestSuite) SetupSuite() {
	transactional := saga.Config(&saga.Opts{
		RabbitUri:    "amqp://rabbit:1234@localhost:5672",
		Microservice: micro.RoomInventory,
		Events: []event.MicroserviceEvent{
			event.SocialNewUserEvent,
			event.SocialBlockChatEvent,
		},
	},
	)
	eventEmitter := transactional.ConnectToEvents()
	suite.e = eventEmitter
	suite.t = transactional
}

func (suite *EventsTestSuite) TestHealthCheck() {
	err := suite.t.HealthCheck()
	suite.Require().NoError(err)
}

func stringPtr(s string) *string {
	return &s
}

func timePtr(t time.Time) *time.Time {
	return &t
}

func (suite *EventsTestSuite) TestSubscribedEvents() {
	testUser := &event.SocialUser{
		ID:              "64f7d1c3b0d8f1b2a4c6a123",
		Username:        "testuser",
		FirstName:       stringPtr("John"),
		LastName:        stringPtr("Doe"),
		Gender:          event.GenderMale,
		IsPublicProfile: true,
		Followers:       []string{"5fc2e9d8474e4c88a5f9d567", "604d22c6741f4e53a8c8e9ba"},
		Following:       []string{"6239f1d547e3f9c8b5e7a91e", "632fc2e3f0b9d1a8c4d3b9f1"},
		Email:           "testuser@example.com",
		Birthday:        timePtr(time.Date(1980, time.January, 1, 0, 0, 0, 0, time.UTC)),
		Location: &event.UserLocation{
			Continent: "North America",
			Country:   "United States",
			Region:    "California",
			City:      "San Francisco",
		},
		Avatar:       stringPtr("https://example.com/avatar.png"),
		SocialMedia:  &event.SocialMedia{"Twitter": "@testuser"},
		Preferences:  []string{"640d91b7e6d3f9c2b4e8a3c5", "65e8c7f4b3a2d5f9c1e0a4b6"},
		BlockedUsers: []string{"00000020f51bb4362eee2a4d"},
		RPMAvatarID:  stringPtr("979dc0f8b1dbf41a0638067f"),
		RPMUserID:    stringPtr("d86ca6512a1142d6afa74193"),
		PaidPriceID:  stringPtr("7302771"),
		CreatedAt:    time.Date(2024, time.November, 7, 14, 18, 25, 723972860, time.UTC),
	}

	eventSocialNewUserReceived := make(chan *event.SocialNewUserPayload)
	eventSocialBlockChatReceived := make(chan *event.SocialBlockChatPayload)

	suite.e.On(event.SocialNewUserEvent, func(handler saga.EventHandler) {
		eventPayload := saga.ParsePayload(handler.Payload, &event.SocialNewUserPayload{})
		defer func(eventPayload *event.SocialNewUserPayload) {
			eventSocialNewUserReceived <- eventPayload
		}(eventPayload)
		handler.Channel.AckMessage()
	})
	suite.e.On(event.SocialBlockChatEvent, func(handler saga.EventHandler) {
		eventPayload := saga.ParsePayload(handler.Payload, &event.SocialBlockChatPayload{})
		defer func(eventPayload *event.SocialBlockChatPayload) {
			eventSocialBlockChatReceived <- eventPayload
		}(eventPayload)
		handler.Channel.AckMessage()
	})

	err := saga.PublishEvent(&event.SocialNewUserPayload{
		SocialUser: *testUser,
	})
	suite.Require().NoError(err)
	err = saga.PublishEvent(&event.SocialBlockChatPayload{
		UserID:        "1234",
		UserToBlockID: "4321",
	})
	suite.Require().NoError(err)

	p1 := <-eventSocialNewUserReceived
	suite.Equal(testUser, &p1.SocialUser)

	p2 := <-eventSocialBlockChatReceived
	suite.Equal("1234", p2.UserID)
	suite.Equal("4321", p2.UserToBlockID)
}

func TestEventsTestSuite(t *testing.T) {
	suite.Run(t, new(EventsTestSuite))
}
