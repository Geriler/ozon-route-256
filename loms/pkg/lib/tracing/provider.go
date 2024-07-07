package tracing

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	sdkResource "go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.25.0"
	"route256/loms/internal/config"
)

func MustLoadTraceProvider(cfg config.Config) *trace.TracerProvider {
	explorer, err := otlptracehttp.New(context.Background(), otlptracehttp.WithEndpointURL(fmt.Sprintf("%s:%d", cfg.Tracer.Host, cfg.Tracer.Port)))
	if err != nil {
		panic(err.Error())
	}

	resource, err := sdkResource.Merge(
		sdkResource.Default(),
		sdkResource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(cfg.ApplicationName),
			semconv.DeploymentEnvironment(cfg.Env),
		),
	)

	traceProvider := trace.NewTracerProvider(
		trace.WithBatcher(explorer),
		trace.WithResource(resource),
	)

	return traceProvider
}
