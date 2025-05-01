package utils

import (
	"context"
	"log"
	"runtime"
	"sample-web/configs"
	"strings"

	"go.opentelemetry.io/otel/trace"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

func InitTracer(config configs.TracingConfig) func(context.Context) error {
	ctx := context.Background()

	otelOptions := []otlptracegrpc.Option{
		otlptracegrpc.WithEndpoint(config.CollectorUrl),
	}

	if config.Insecure {
		otelOptions = append(otelOptions, otlptracegrpc.WithInsecure())
	}

	exporter, err := otlptracegrpc.New(ctx, otelOptions...)

	if err != nil {
		log.Fatalf("failed to create exporter: %v", err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(config.ServiceName),
		)),
	)

	otel.SetTracerProvider(tp)

	return tp.Shutdown
}

func Tracer() trace.Tracer {
	pc, _, _, ok := runtime.Caller(1)
	if !ok {
		return otel.Tracer("unknown")
	}

	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return otel.Tracer("unknown")
	}

	fullFuncName := fn.Name()
	parts := strings.Split(fullFuncName, "/")
	if len(parts) == 0 {
		return otel.Tracer("unknown")
	}

	pkgPath := strings.Join(parts[len(parts)-2:], "/")
	pkgOnly := strings.Split(pkgPath, ".")[0]

	return otel.Tracer(pkgOnly)
}
