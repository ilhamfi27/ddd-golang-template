package config

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
)

// InitTracing initializes the OpenTelemetry tracer with Jaeger exporter
func InitTracing(serviceName, jaegerEndpoint string) (func(context.Context) error, error) {
	// Create OTLP exporter that sends to Jaeger
	// WithEndpoint expects host:port (no scheme), WithInsecure() enables plain HTTP
	exporter, err := otlptracehttp.New(context.Background(),
		otlptracehttp.WithEndpoint(jaegerEndpoint),
		otlptracehttp.WithInsecure(),
	)
	if err != nil {
		return nil, err
	}

	// Create resource with service metadata
	res, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			"https://opentelemetry.io/schemas/1.26.0",
			attribute.String("service.name", serviceName),
			attribute.String("service.version", "1.0.0"),
		),
	)
	if err != nil {
		return nil, err
	}

	// Create batch span processor for efficient batching
	batcher := trace.NewBatchSpanProcessor(exporter)

	// Create trace provider
	tp := trace.NewTracerProvider(
		trace.WithResource(res),
		trace.WithSpanProcessor(batcher),
		// Sample all traces in development, adjust for production
		trace.WithSampler(trace.AlwaysSample()),
	)

	// Set as global tracer provider
	otel.SetTracerProvider(tp)

	// Return shutdown function for graceful shutdown
	return tp.Shutdown, nil
}
