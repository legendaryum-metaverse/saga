package test

import (
	"context"
	"testing"

	"github.com/legendaryum-metaverse/saga"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func recordingInvoker(recorded *string) grpc.UnaryInvoker {
	return func(ctx context.Context, _ string, _, _ interface{}, _ *grpc.ClientConn, _ ...grpc.CallOption) error {
		if md, ok := metadata.FromOutgoingContext(ctx); ok {
			if values := md.Get(saga.OperationHeader); len(values) > 0 {
				*recorded = values[0]
			}
		}
		return nil
	}
}

func TestUnaryClientOperationInterceptorAttachesTheHeader(t *testing.T) {
	ctx := saga.WithOperation(context.Background(), "op-123")
	var recorded string

	err := saga.UnaryClientOperationInterceptor(
		ctx, "/Service/Method", nil, nil, nil, recordingInvoker(&recorded),
	)

	require.NoError(t, err)
	assert.Equal(t, "op-123", recorded)
}

func TestUnaryClientOperationInterceptorLeavesMetadataUntouchedWithoutOperation(t *testing.T) {
	var recorded string

	err := saga.UnaryClientOperationInterceptor(
		context.Background(), "/Service/Method", nil, nil, nil, recordingInvoker(&recorded),
	)

	require.NoError(t, err)
	assert.Empty(t, recorded)
}

func echoingHandler(seen *string) grpc.UnaryHandler {
	return func(ctx context.Context, _ interface{}) (interface{}, error) {
		if operationID, ok := saga.OperationFrom(ctx); ok {
			*seen = operationID
		}
		return nil, nil
	}
}

func TestUnaryServerOperationInterceptorBindsFromIncomingMetadata(t *testing.T) {
	ctx := metadata.NewIncomingContext(
		context.Background(),
		metadata.Pairs(saga.OperationHeader, "op-from-metadata"),
	)
	var seen string

	_, err := saga.UnaryServerOperationInterceptor(
		ctx, nil, &grpc.UnaryServerInfo{}, echoingHandler(&seen),
	)

	require.NoError(t, err)
	assert.Equal(t, "op-from-metadata", seen)
}

func TestUnaryServerOperationInterceptorWithoutMetadataBindsNothing(t *testing.T) {
	var seen string

	_, err := saga.UnaryServerOperationInterceptor(
		context.Background(), nil, &grpc.UnaryServerInfo{}, echoingHandler(&seen),
	)

	require.NoError(t, err)
	assert.Empty(t, seen)
}

// The chain that matters end to end: a server receives an operation, and a
// further outgoing call made from inside that handler carries it onward
// without the handler doing anything special — mirrors the equivalent test
// in rust-library's legend-saga::grpc.
func TestServerInterceptorThenClientInterceptorForwardsAcrossAHop(t *testing.T) {
	var downstreamSaw string

	handler := func(ctx context.Context, _ interface{}) (interface{}, error) {
		return nil, saga.UnaryClientOperationInterceptor(
			ctx, "/Downstream/Method", nil, nil, nil, recordingInvoker(&downstreamSaw),
		)
	}

	incomingCtx := metadata.NewIncomingContext(
		context.Background(),
		metadata.Pairs(saga.OperationHeader, "op-hopping-through"),
	)

	_, err := saga.UnaryServerOperationInterceptor(
		incomingCtx, nil, &grpc.UnaryServerInfo{}, handler,
	)

	require.NoError(t, err)
	assert.Equal(t, "op-hopping-through", downstreamSaw)
}
