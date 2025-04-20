package middlewares

import (
	"context"
	"fmt"
	"log"

	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func enrichContextWithTraceID(ctx context.Context) (context.Context, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ctx, fmt.Errorf("no metadata")
	}

	traceIDString := md["x-trace-id"][0]
	// Convert string to byte array
	traceID, err := trace.TraceIDFromHex(traceIDString)
	if err != nil {
		return ctx, fmt.Errorf("invalid trace ID")
	}

	// Creating a span context with a predefined trace-id
	spanContext := trace.NewSpanContext(trace.SpanContextConfig{
		TraceID: traceID,
	})
	// Embedding span config into the context
	ctx = trace.ContextWithSpanContext(ctx, spanContext)

	log.Printf("find traceID: %v", traceIDString)
	return ctx, nil
}

func Tracing(ctx context.Context, req any, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	ctx, err = enrichContextWithTraceID(ctx)
	if err != nil {
		log.Printf("enrichContextWithTraceID: %v", err)
	}

	return handler(ctx, req)
}
