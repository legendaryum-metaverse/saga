package saga

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// UnaryClientOperationInterceptor attaches the operation bound to ctx as
// outgoing gRPC metadata under OperationHeader, mirroring the AMQP header
// mechanism on a second transport. A call made without an operation in ctx
// is left untouched — same permissive-migration stance as applyOperationHeader.
//
// Wire it once per client connection:
//
//	grpc.NewClient(addr, grpc.WithChainUnaryInterceptor(saga.UnaryClientOperationInterceptor))
func UnaryClientOperationInterceptor(
	ctx context.Context,
	method string,
	req, reply interface{},
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption,
) error {
	if operationID, ok := OperationFrom(ctx); ok {
		ctx = metadata.AppendToOutgoingContext(ctx, OperationHeader, operationID)
	}
	return invoker(ctx, method, req, reply, cc, opts...)
}

// UnaryServerOperationInterceptor reads OperationHeader off the incoming call
// and binds it to the context the handler receives, so any code inside the
// handler — including a further outgoing gRPC call — sees it via
// OperationFrom without touching metadata itself.
//
// Wire it once per server, before any interceptor whose logic depends on the
// operation being present:
//
//	grpc.NewServer(grpc.ChainUnaryInterceptor(saga.UnaryServerOperationInterceptor, ...))
func UnaryServerOperationInterceptor(
	ctx context.Context,
	req interface{},
	_ *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if values := md.Get(OperationHeader); len(values) > 0 && values[0] != "" {
			ctx = WithOperation(ctx, values[0])
		}
	}
	return handler(ctx, req)
}
