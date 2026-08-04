package saga

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
)

// OperationHeader is the AMQP header carrying the SIPLEI operation identifier.
// The same key is used for gRPC metadata so there is a single name across transports.
const OperationHeader = "x-operation-id"

type operationContextKey struct{}

// MissingOperationHook is invoked when a consumed message carries no operation.
// It is a hook instead of a direct metrics dependency so the library stays free
// of a metrics backend. Services wire it to `events_without_operation_total`.
type MissingOperationHook func(microservice, eventType string)

var missingOperationHook MissingOperationHook

// SetMissingOperationHook registers the callback for messages without an operation.
// Not safe for concurrent use with message consumption; call it during setup.
func SetMissingOperationHook(hook MissingOperationHook) {
	missingOperationHook = hook
}

func reportMissingOperation(microservice, eventType string) {
	if missingOperationHook != nil {
		missingOperationHook(microservice, eventType)
	}
}

// WithOperation returns a context carrying the operation identifier.
func WithOperation(ctx context.Context, operationID string) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithValue(ctx, operationContextKey{}, operationID)
}

// OperationFrom extracts the operation identifier from a context.
func OperationFrom(ctx context.Context) (string, bool) {
	if ctx == nil {
		return "", false
	}
	operationID, ok := ctx.Value(operationContextKey{}).(string)
	if !ok || operationID == "" {
		return "", false
	}
	return operationID, true
}

// operationFromHeaders reads the operation identifier off an AMQP message.
func operationFromHeaders(headers amqp.Table) string {
	if headers == nil {
		return ""
	}
	operationID, ok := headers[OperationHeader].(string)
	if !ok {
		return ""
	}
	return operationID
}

// detachOperation returns a background context carrying only the operation.
// Used before spawning goroutines that outlive the caller's context, so the
// operation survives without inheriting its cancellation.
func detachOperation(ctx context.Context) context.Context {
	operationID, ok := OperationFrom(ctx)
	if !ok {
		return context.Background()
	}
	return WithOperation(context.Background(), operationID)
}

// applyOperationHeader adds the operation header when the context carries one.
// Messages published without an operation are left untouched, which keeps the
// permissive migration window working.
func applyOperationHeader(ctx context.Context, headers amqp.Table) amqp.Table {
	operationID, ok := OperationFrom(ctx)
	if !ok {
		return headers
	}
	if headers == nil {
		headers = amqp.Table{}
	}
	headers[OperationHeader] = operationID
	return headers
}
