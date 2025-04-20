package tracing

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	otelTrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"

	"route256/cart/internal/logger"
)

func init() {
	ctx := context.Background()
	logger.Warnw(ctx, "init tracing")

	jaegerResource, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName("cart"),
		),
	)
	if err != nil {
		panic(err)
	}

	exp, err := otlptracehttp.New(
		ctx,
		otlptracehttp.WithEndpointURL("http://localhost:4318"), // это можно вынести в конфиг.
	)
	if err != nil {
		panic(err)
	}

	traceProvider := otelTrace.NewTracerProvider(
		otelTrace.WithBatcher(exp),
		otelTrace.WithResource(jaegerResource),
	)

	otel.SetTracerProvider(traceProvider)
	//Важно! Без MapPropagator не будет работать.
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
}

func StartFromContext(ctx context.Context, spanName string) (context.Context, trace.Span) {
	return otel.GetTracerProvider().Tracer("kek").Start(ctx, spanName)
}
