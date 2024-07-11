package tracing

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	sdkResource "go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.25.0"
	"route256/cart/internal/config"
)

func MustLoadTraceProvider(cfg config.Config) *trace.TracerProvider {
	url := fmt.Sprintf("%s:%d", cfg.Tracer.Host, cfg.Tracer.Port)

	explorer, err := otlptracehttp.New(
		context.Background(),
		otlptracehttp.WithEndpoint(url),
		otlptracehttp.WithInsecure(),
	)
	if err != nil {
		panic(err.Error())
	}

	otel.SetTextMapPropagator(propagation.TraceContext{})

	traceProvider := trace.NewTracerProvider(
		trace.WithBatcher(explorer),
		trace.WithResource(sdkResource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(cfg.ApplicationName),
			semconv.DeploymentEnvironment(cfg.Env),
		)),
	)

	otel.SetTracerProvider(traceProvider)

	return traceProvider
}
