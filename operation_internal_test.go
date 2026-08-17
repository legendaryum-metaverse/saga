package saga

import (
	"context"
	"testing"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOperationFromHeaders(t *testing.T) {
	tests := []struct {
		name    string
		headers amqp.Table
		want    string
	}{
		{name: "header present", headers: amqp.Table{OperationHeader: "op-123"}, want: "op-123"},
		{name: "header absent", headers: amqp.Table{"all-micro": "yes"}, want: ""},
		{name: "nil headers", headers: nil, want: ""},
		{name: "non string value", headers: amqp.Table{OperationHeader: 42}, want: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, operationFromHeaders(tt.headers))
		})
	}
}

func TestApplyOperationHeaderWithOperation(t *testing.T) {
	ctx := WithOperation(context.Background(), "op-123")
	headers := amqp.Table{"all-micro": "yes"}

	got := applyOperationHeader(ctx, headers)

	assert.Equal(t, "op-123", got[OperationHeader])
	assert.Equal(t, "yes", got["all-micro"])
}

func TestApplyOperationHeaderWithoutOperationLeavesHeadersUntouched(t *testing.T) {
	headers := amqp.Table{"all-micro": "yes"}

	got := applyOperationHeader(context.Background(), headers)

	assert.NotContains(t, got, OperationHeader)
	assert.Len(t, got, 1)
}

func TestApplyOperationHeaderOnNilTable(t *testing.T) {
	ctx := WithOperation(context.Background(), "op-123")

	got := applyOperationHeader(ctx, nil)

	assert.Equal(t, "op-123", got[OperationHeader])
}

func TestApplyOperationHeaderOnNilTableWithoutOperationStaysNil(t *testing.T) {
	assert.Nil(t, applyOperationHeader(context.Background(), nil))
}

func TestOperationHeaderIsNotMistakenForAnEventKey(t *testing.T) {
	// findEventValues scans every string header looking for known events.
	// The operation header must never be picked up as one.
	headers := amqp.Table{
		OperationHeader: "op-123",
		"AUTH.NEW_USER": "auth.new_user",
	}

	events, err := findEventValues(headers)

	require.NoError(t, err)
	assert.Len(t, events, 1)
}

func TestDetachOperation(t *testing.T) {
	parent, cancel := context.WithCancel(WithOperation(context.Background(), "op-123"))
	detached := detachOperation(parent)
	cancel()

	operationID, ok := OperationFrom(detached)

	assert.True(t, ok)
	assert.Equal(t, "op-123", operationID)
	require.NoError(t, detached.Err(), "detached context must survive parent cancellation")
}

func TestMissingOperationHook(t *testing.T) {
	t.Cleanup(func() { SetMissingOperationHook(nil) })

	var gotMicroservice, gotEvent string
	calls := 0
	SetMissingOperationHook(func(microservice, eventType string) {
		calls++
		gotMicroservice = microservice
		gotEvent = eventType
	})

	reportMissingOperation("social", "auth.new_user")

	assert.Equal(t, 1, calls)
	assert.Equal(t, "social", gotMicroservice)
	assert.Equal(t, "auth.new_user", gotEvent)
}

func TestMissingOperationHookUnset(t *testing.T) {
	SetMissingOperationHook(nil)

	assert.NotPanics(t, func() { reportMissingOperation("social", "auth.new_user") })
}
