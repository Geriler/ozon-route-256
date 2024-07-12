package middleware

import (
	"context"
	"log"

	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"route256/loms/pkg/lib/tracing"
)

func Tracing(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	ctx = getCtxWithTraceID(ctx)
	ctx, span := tracing.StartSpanFromContext(ctx, info.FullMethod, trace.WithSpanKind(trace.SpanKindServer))
	defer span.End()

	return handler(ctx, req)
}

func getCtxWithTraceID(ctx context.Context) context.Context {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ctx
	}
	traceIdStrings := md.Get("x-trace-id")
	if len(traceIdStrings) == 0 {
		return ctx
	}
	traceId, err := trace.TraceIDFromHex(traceIdStrings[0])
	log.Println(traceId)
	if err != nil {
		return ctx
	}

	spanContext := trace.NewSpanContext(trace.SpanContextConfig{
		TraceID: traceId,
	})

	ctx = trace.ContextWithSpanContext(ctx, spanContext)

	return ctx
}
