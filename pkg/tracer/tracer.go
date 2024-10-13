package tracer

import (
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"gopkg.in/natefinch/lumberjack.v2"
)

// InitTracer initializes the tracer provider and returns it.
func InitTracer() (*sdktrace.TracerProvider, error) {
	rotatingLogFile := &lumberjack.Logger{
		Filename:   "trace.log",
		MaxSize:    20,
		MaxBackups: 5,
		MaxAge:     30,
		Compress:   true,
		LocalTime:  true,
	}

	exporter, err := stdouttrace.New(stdouttrace.WithWriter(rotatingLogFile), stdouttrace.WithPrettyPrint())
	if err != nil {
		return nil, err
	}

	sampler := sdktrace.ParentBased(sdktrace.TraceIDRatioBased(0.05))

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithSampler(sampler),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceInstanceIDKey.String(uuid.NewString()),
			semconv.ServiceNameKey.String("zenith-server"),
		)),
	)

	otel.SetTracerProvider(tp)
	return tp, nil
}
