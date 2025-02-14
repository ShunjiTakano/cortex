package cortex

import (
	"context"

	"github.com/opentracing/opentracing-go"
	objstoretracing "github.com/thanos-io/objstore/tracing"
	"github.com/thanos-io/thanos/pkg/tracing"
	"google.golang.org/grpc"
)

// ThanosTracerUnaryInterceptor injects the opentracing global tracer into the context
// in order to get it picked up by Thanos components.
func ThanosTracerUnaryInterceptor(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	ctx = objstoretracing.ContextWithTracer(ctx, opentracing.GlobalTracer())
	return handler(tracing.ContextWithTracer(ctx, opentracing.GlobalTracer()), req)
}

// ThanosTracerStreamInterceptor injects the opentracing global tracer into the context
// in order to get it picked up by Thanos components.
func ThanosTracerStreamInterceptor(srv interface{}, ss grpc.ServerStream, _ *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	ctx := objstoretracing.ContextWithTracer(ss.Context(), opentracing.GlobalTracer())
	return handler(srv, wrappedServerStream{
		ctx:          tracing.ContextWithTracer(ctx, opentracing.GlobalTracer()),
		ServerStream: ss,
	})
}

type wrappedServerStream struct {
	ctx context.Context
	grpc.ServerStream
}

func (ss wrappedServerStream) Context() context.Context {
	return ss.ctx
}
