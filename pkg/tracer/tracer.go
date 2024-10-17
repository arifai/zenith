package tracer

import (
	"github.com/arifai/zenith/config"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"time"
)

// InitTracer initializes the tracer provider and returns it.
func InitTracer(config *config.Config) (*sdktrace.TracerProvider, error) {

	exporter, err := zipkin.New(config.ZipkinURL)
	if err != nil {
		return nil, err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(
			exporter,
			sdktrace.WithMaxExportBatchSize(sdktrace.DefaultMaxExportBatchSize),
			sdktrace.WithBatchTimeout(sdktrace.DefaultScheduleDelay*time.Millisecond),
		),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceInstanceIDKey.String(uuid.NewString()),
			semconv.ServiceNameKey.String("zenith-server"),
		)),
	)

	otel.SetTracerProvider(tp)
	return tp, nil
}
