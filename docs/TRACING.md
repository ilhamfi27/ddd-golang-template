# Distributed Tracing Guide

## Overview

Distributed tracing helps you understand request flows across your application and services. It answers questions like:
- How long did each operation take?
- Where did the request fail?
- Which services did the request go through?
- What was the dependency chain?

## Architecture: OpenTelemetry + Jaeger

```
┌─────────────────────────────────────────────────────────┐
│                    Your Go Application                   │
│  (with OpenTelemetry instrumentation)                   │
└──────────────────┬──────────────────────────────────────┘
                   │
                   │ Sends trace data
                   │ (spans, metrics, logs)
                   ↓
        ┌──────────────────────┐
        │  Jaeger Collector    │  (Docker Container)
        │  - OTLP HTTP:4318    │
        │  - UI:16686          │
        └──────────────────────┘
                   │
                   │ Stores & processes
                   ↓
        ┌──────────────────────┐
        │  Jaeger Storage      │  (In-memory or database)
        └──────────────────────┘
                   │
                   │ Query & visualize
                   ↓
    ┌──────────────────────────────┐
    │  Jaeger Web UI (localhost:16686)
    │  - Search traces
    │  - View service dependencies
    │  - Analyze performance
    └──────────────────────────────┘
```

## Components Explained

### 1. **OpenTelemetry** (Code Library)
- **What it is**: A set of APIs and SDKs for instrumenting your code
- **What it does**: Captures trace data from your application
- **Language**: Go, Python, JavaScript, etc.
- **Package**: `go.opentelemetry.io/otel/*`

```go
// Example: OpenTelemetry in your code
ctx, span := tracer.Start(context.Background(), "operation-name")
defer span.End()

span.SetAttributes(
    trace.String("user_id", "123"),
    trace.Int("items_count", 5),
)
```

### 2. **Jaeger** (Tracing Backend)
- **What it is**: A storage and visualization system for traces
- **What it does**: 
  - Collects trace data from your application
  - Stores it
  - Provides a web UI to view traces
- **Deployment**: Usually runs as a Docker container
- **UI**: Available at `http://localhost:16686`

### 3. **OTLP** (Protocol)
- **Full name**: OpenTelemetry Protocol
- **What it is**: The protocol that sends data from your app to Jaeger
- **Port**: 4318 (HTTP) or 4317 (gRPC)
- **Why**: Vendor-neutral, works with any backend

## Installation Steps

### Step 1: Add Go Dependencies
```bash
go get go.opentelemetry.io/otel
go get go.opentelemetry.io/otel/sdk/trace
go get go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp
go get go.opentelemetry.io/otel/sdk/resource
go get go.opentelemetry.io/otel/exporters/jaeger/otlphttp
```

### Step 2: Run Jaeger with Docker

**Option A: Docker Run (Quick)**
```bash
docker run -d \
  -p 16686:16686 \
  -p 4318:4318 \
  jaegertracing/all-in-one
```

**Option B: Docker Compose (Recommended)**
```yaml
# docker-compose.yml
version: '3'
services:
  jaeger:
    image: jaegertracing/all-in-one
    ports:
      - "16686:16686"  # Web UI
      - "4318:4318"    # OTLP HTTP receiver
      - "6831:6831/udp" # Jaeger compact Thrift
    environment:
      - COLLECTOR_OTLP_ENABLED=true
```

Run: `docker-compose up -d jaeger`

### Step 3: Instrument Your Application (Code)
```go
// internal/config/tracing.go
package config

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlphttp"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/semconv/v1.10.0"
)

func InitTracing(serviceName string) (func(context.Context) error, error) {
	// Create OTLP exporter (sends to Jaeger)
	exporter, err := otlphttp.New(context.Background(),
		otlphttp.WithEndpoint("http://localhost:4318"),
	)
	if err != nil {
		return nil, err
	}

	// Create resource (service metadata)
	resource, err := resource.Merge(
		resource.Default(),
		resource.NewSchemaURLResource(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(serviceName),
		),
	)
	if err != nil {
		return nil, err
	}

	// Create trace provider
	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(resource),
	)

	// Set as global
	otel.SetTracerProvider(tp)

	// Return shutdown function
	return tp.Shutdown, nil
}
```

### Step 4: Use Tracing in Your Middleware
```go
// internal/handlers/middleware/tracing.go
package middleware

import (
	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

func TracingMiddleware() echo.MiddlewareFunc {
	tracer := otel.Tracer("ddd-golang-template")
	
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Create a span for this request
			ctx, span := tracer.Start(c.Request().Context(), c.Request().URL.Path)
			defer span.End()

			// Add attributes to the span
			span.SetAttributes(
				trace.String("http.method", c.Request().Method),
				trace.String("http.url", c.Request().URL.String()),
				trace.String("http.client_ip", c.RealIP()),
			)

			// Continue with request
			c.SetRequest(c.Request().WithContext(ctx))
			err := next(c)

			// Record response status
			if err == nil {
				span.SetAttributes(trace.Int("http.status_code", c.Response().Status))
			} else {
				span.RecordError(err)
			}

			return err
		}
	}
}
```

### Step 5: Initialize in Main
```go
// cmd/main.go
package main

import (
	"context"
	"github.com/ilhamfi27/ddd-golang-template/internal/config"
	"github.com/ilhamfi27/ddd-golang-template/internal/handlers/middleware"
)

func main() {
	// Initialize tracing
	shutdown, err := config.InitTracing("ddd-golang-template")
	if err != nil {
		panic(err)
	}
	defer shutdown(context.Background())

	// ... rest of your application
	h := handlers.NewHandler()
	
	// Add tracing middleware
	h.app.Use(middleware.TracingMiddleware())
	
	h.Start()
}
```

## How It Works (Step by Step)

### When a Request Comes In:

```
1. HTTP Request arrives
   ↓
2. Middleware creates a "span" (trace event)
   └─ Span name: "/api/users"
   └─ Attributes: method, URL, client IP
   ↓
3. Your code executes
   └─ If you create child spans, they're nested
   ↓
4. Response is sent
   └─ Span records status code
   └─ Span records duration
   ↓
5. Middleware ends the span
   ↓
6. OpenTelemetry batches the span data
   ↓
7. Data is sent to Jaeger via OTLP protocol (HTTP on port 4318)
   ↓
8. Jaeger stores it
   ↓
9. You view it in UI at http://localhost:16686
```

## Creating Child Spans (Nested Operations)

```go
// Example: Tracing database operations
func (s *UserService) CreateUser(ctx context.Context, user *User) error {
	tracer := otel.Tracer("user-service")
	
	// Create child span for database operation
	ctx, span := tracer.Start(ctx, "database.create_user")
	defer span.End()
	
	err := s.repo.Create(ctx, user)
	if err != nil {
		span.RecordError(err)
		return err
	}
	
	return nil
}
```

## Viewing Traces

### Open Jaeger UI
```
http://localhost:16686
```

### Steps:
1. **Select Service**: Choose "ddd-golang-template" from dropdown
2. **View Traces**: All traces from your app appear here
3. **Click a Trace**: See the full request flow
4. **Analyze**: See which operations took longest

### What You'll See:

```
Trace: POST /api/users (completed in 250ms)
├─ POST /api/users (100ms)
│  ├─ database.create_user (45ms)
│  └─ email.send (150ms)
└─ HTTP response (50ms)
```

## Common Ports

| Port | Service | Protocol | Purpose |
|------|---------|----------|---------|
| 4318 | Jaeger | HTTP | OTLP data from your app |
| 4317 | Jaeger | gRPC | OTLP data (alternative) |
| 16686 | Jaeger | HTTP | Web UI |
| 6831 | Jaeger | UDP | Legacy Thrift protocol |

## Environment Variables

```env
# Jaeger endpoint
OTEL_EXPORTER_OTLP_ENDPOINT=http://localhost:4318

# Service name
OTEL_SERVICE_NAME=ddd-golang-template

# Sampling rate (0-1, where 1 = trace everything)
OTEL_TRACES_SAMPLER=always_on

# Or sample only 10% of traces
OTEL_TRACES_SAMPLER=traceidratio
OTEL_TRACES_SAMPLER_ARG=0.1
```

## Architecture in Your Codebase

```
internal/
├── config/
│   ├── tracing.go          ← Initialize tracer
│   └── app/
│       └── env.go          ← Load env variables
├── handlers/
│   ├── middleware/
│   │   └── tracing.go      ← Tracing middleware
│   └── rest.go             ← Register middleware
└── ...
```

## Production Considerations

### 1. **Sampling**
Don't trace 100% of requests in production (too much data):
```go
// Sample only 10% of requests
sampler := sdktrace.NewParentBasedSampler(
	sdktrace.TraceIDRatioBased(0.1),
)

tp := sdktrace.NewTracerProvider(
	sdktrace.WithSampler(sampler),
	// ...
)
```

### 2. **Batch Processing**
Use batch exporter (not default):
```go
batcher := sdktrace.NewBatchSpanProcessor(exporter)
tp := sdktrace.NewTracerProvider(
	sdktrace.WithSpanProcessor(batcher),
	// ...
)
```

### 3. **Production Jaeger**
Use a proper Jaeger deployment:
- **Elasticsearch** for storage (instead of in-memory)
- **Kafka** for scalability
- **Cassandra** for persistence
- **Load balancer** for multiple Jaeger instances

### 4. **Log Correlation**
Link logs and traces with trace ID:
```go
// In your logging
traceID := trace.SpanContextFromContext(ctx).TraceID()
logger.Info("operation completed", 
	slog.String("trace_id", traceID.String()),
)
```

## Troubleshooting

### Traces not appearing in Jaeger?

1. **Check Docker is running**
   ```bash
   docker ps | grep jaeger
   ```

2. **Check connectivity**
   ```bash
   curl http://localhost:16686
   ```

3. **Check service name**
   - In Jaeger UI, select dropdown - should see your service name
   - If empty, service name might be misconfigured

4. **Enable debug logging**
   ```go
   import "log/slog"
   handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
   slog.SetDefault(slog.New(handler))
   ```

### High memory usage?
- Reduce sampling rate
- Use batch span processor
- Reduce span buffer size

## Next Steps

1. **Install dependencies**: `go get go.opentelemetry.io/otel/...`
2. **Run Jaeger**: `docker-compose up -d jaeger`
3. **Add tracing config**: Create `internal/config/tracing.go`
4. **Add middleware**: Create `internal/handlers/middleware/tracing.go`
5. **Initialize in main**: Update `cmd/main.go`
6. **Test**: Make requests and view traces at http://localhost:16686

## Further Reading

- [OpenTelemetry Documentation](https://opentelemetry.io/docs/instrumentation/go/)
- [Jaeger Documentation](https://www.jaegertracing.io/docs/)
- [OTEL Best Practices](https://opentelemetry.io/docs/reference/specification/protocol/exporter/)

---

**Key Takeaway**: 
- **OpenTelemetry** = The code library (instruments your app)
- **Jaeger** = The backend (collects & visualizes traces)
- **OTLP** = The protocol that connects them
- **UI** = http://localhost:16686 (where you see the magic happen)
