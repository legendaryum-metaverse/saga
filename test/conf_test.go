package test

import (
	"testing"

	"github.com/legendaryum-metaverse/saga"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"github.com/legendaryum-metaverse/saga/event"
	"github.com/legendaryum-metaverse/saga/micro"
)

func TestConfig(t *testing.T) {
	tests := []struct {
		name         string
		opts         *saga.Opts
		shouldPanic  bool
		panicMessage string
	}{
		{
			name: "Valid options",
			opts: &saga.Opts{
				RabbitUri:    "amqp://guest:guest@localhost:5672/",
				Microservice: micro.Auth,
				Events:       []event.MicroserviceEvent{},
				Logger:       zap.NewNop(),
			},
			shouldPanic: false,
		},
		{
			name: "Valid options without Events",
			opts: &saga.Opts{
				RabbitUri:    "amqp://guest:guest@localhost:5672/",
				Microservice: micro.Auth,
				Logger:       zap.NewNop(),
			},
			shouldPanic: false,
		},
		{
			name: "Missing RabbitUri",
			opts: &saga.Opts{
				Microservice: micro.Auth,
				Events:       []event.MicroserviceEvent{},
				Logger:       zap.NewNop(),
			},
			shouldPanic:  true,
			panicMessage: "Invalid options: Key: 'Opts.RabbitUri' Error:Field validation for 'RabbitUri' failed on the 'required' tag",
		},
		{
			name: "Invalid RabbitUri",
			opts: &saga.Opts{
				RabbitUri:    "invalid-uri",
				Microservice: micro.Auth,
				Events:       []event.MicroserviceEvent{},
				Logger:       zap.NewNop(),
			},
			shouldPanic:  true,
			panicMessage: "Invalid options: Key: 'Opts.RabbitUri' Error:Field validation for 'RabbitUri' failed on the 'url' tag",
		},
		{
			name: "Missing Microservice",
			opts: &saga.Opts{
				RabbitUri: "amqp://guest:guest@localhost:5672/",
				Events:    []event.MicroserviceEvent{},
				Logger:    zap.NewNop(),
			},
			shouldPanic:  true,
			panicMessage: "Invalid options: Key: 'Opts.Microservice' Error:Field validation for 'Microservice' failed on the 'required' tag",
		},
		{
			name: "Invalid Microservice",
			opts: &saga.Opts{
				RabbitUri:    "amqp://guest:guest@localhost:5672/",
				Microservice: micro.AvailableMicroservices("invalid-service"),
				Events:       []event.MicroserviceEvent{},
				Logger:       zap.NewNop(),
			},
			shouldPanic:  true,
			panicMessage: "Invalid options: Key: 'Opts.Microservice' Error:Field validation for 'Microservice' failed on the 'microservice' tag",
		},
		{
			name: "Missing Logger",
			opts: &saga.Opts{
				RabbitUri:    "amqp://guest:guest@localhost:5672/",
				Microservice: micro.Auth,
				Events:       []event.MicroserviceEvent{},
			},
			shouldPanic:  true,
			panicMessage: "Invalid options: Key: 'Opts.Logger' Error:Field validation for 'Logger' failed on the 'required' tag",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.shouldPanic {
				assert.PanicsWithValue(t, tt.panicMessage, func() { saga.Config(tt.opts) })
			} else {
				assert.NotPanics(t, func() {
					tr := saga.Config(tt.opts)
					assert.NotNil(t, tr)
					assert.Equal(t, tt.opts.RabbitUri, tr.RabbitUri)
					assert.Equal(t, tt.opts.Microservice, tr.Microservice)
					assert.Equal(t, tt.opts.Events, tr.Events)
					assert.Equal(t, tt.opts.Logger, tr.Logger)
				})
			}
		})
	}
}
