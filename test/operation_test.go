package test

import (
	"context"
	"testing"

	"github.com/legendaryum-metaverse/saga"
	"github.com/stretchr/testify/assert"
)

func TestOperationContextRoundTrip(t *testing.T) {
	ctx := saga.WithOperation(context.Background(), "op-123")

	operationID, ok := saga.OperationFrom(ctx)

	assert.True(t, ok)
	assert.Equal(t, "op-123", operationID)
}

func TestOperationFromWithoutOperation(t *testing.T) {
	tests := []struct {
		name string
		ctx  context.Context
	}{
		{name: "background context", ctx: context.Background()},
		{name: "nil context", ctx: nil},
		{name: "empty operation", ctx: saga.WithOperation(context.Background(), "")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			operationID, ok := saga.OperationFrom(tt.ctx)

			assert.False(t, ok)
			assert.Empty(t, operationID)
		})
	}
}

func TestWithOperationOnNilContext(t *testing.T) {
	//nolint:staticcheck // exercising the nil-context guard on purpose
	ctx := saga.WithOperation(nil, "op-456")

	operationID, ok := saga.OperationFrom(ctx)

	assert.True(t, ok)
	assert.Equal(t, "op-456", operationID)
}

func TestOperationOverride(t *testing.T) {
	ctx := saga.WithOperation(context.Background(), "op-first")
	ctx = saga.WithOperation(ctx, "op-second")

	operationID, _ := saga.OperationFrom(ctx)

	assert.Equal(t, "op-second", operationID)
}

func TestOperationHeaderName(t *testing.T) {
	// The same key is used for AMQP headers and gRPC metadata across the three
	// libraries. A rename here breaks cross-language propagation.
	assert.Equal(t, "x-operation-id", saga.OperationHeader)
}
