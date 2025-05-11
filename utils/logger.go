package utils

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sample-web/configs"
	"sync"

	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
)

type Logger interface {
	Debug(ctx context.Context, args ...interface{})
	Info(ctx context.Context, args ...interface{})
	Warn(ctx context.Context, args ...interface{})
	Error(ctx context.Context, args ...interface{})
	Fatal(ctx context.Context, args ...interface{})
	Tracer() trace.Tracer
}

type logrusWrapper struct {
	logger *logrus.Logger
	tracer trace.Tracer
}

var (
	instance Logger
	once     sync.Once
)

func InitLogger(config configs.TracingConfig) {
	once.Do(func() {
		otelOptions := []otlptracegrpc.Option{
			otlptracegrpc.WithEndpoint(config.CollectorUrl),
		}

		if config.Insecure {
			otelOptions = append(otelOptions, otlptracegrpc.WithInsecure())
		}
		exporter, err := otlptracegrpc.New(context.Background(), otelOptions...)
		if err != nil {
			panic("failed to create exporter: " + err.Error())
		}

		tp := sdktrace.NewTracerProvider(
			sdktrace.WithBatcher(exporter),
			sdktrace.WithResource(resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceName(config.ServiceName),
			)),
		)

		otel.SetTracerProvider(tp)
		logger := logrus.New()

		instance = &logrusWrapper{
			logger: logger,
			tracer: otel.Tracer(config.ServiceName),
		}
	})
}

func GetLogger() Logger {
	if instance == nil {
		panic("Logger not initialized. Call InitLogger(serviceName) first.")
	}
	return instance
}

func (l *logrusWrapper) logWithTrace(ctx context.Context, level logrus.Level, args ...interface{}) {
	span := trace.SpanFromContext(ctx)
	msg := fmt.Sprint(args...)

	if span.SpanContext().IsValid() {
		// Attach trace_id and span_id to log fields
		fields := logrus.Fields{
			"trace_id": span.SpanContext().TraceID().String(),
			"span_id":  span.SpanContext().SpanID().String(),
		}
		l.logger.WithFields(fields).Log(level, msg)

		// Record error if log level is error or higher
		if level == logrus.ErrorLevel {
			span.RecordError(errors.New(msg))
		} else {
			// Add message as event to span
			span.AddEvent("log", trace.WithAttributes(
				attribute.String("log.severity", level.String()),
				attribute.String("log.message", msg),
			))
		}
		
	} else {
		l.logger.Log(level, msg)
	}
}

func (l *logrusWrapper) Debug(ctx context.Context, args ...interface{}) {
	l.logWithTrace(ctx, logrus.DebugLevel, args...)
}

func (l *logrusWrapper) Info(ctx context.Context, args ...interface{}) {
	l.logWithTrace(ctx, logrus.InfoLevel, args...)
}

func (l *logrusWrapper) Warn(ctx context.Context, args ...interface{}) {
	l.logWithTrace(ctx, logrus.WarnLevel, args...)
}

func (l *logrusWrapper) Error(ctx context.Context, args ...interface{}) {
	l.logWithTrace(ctx, logrus.ErrorLevel, args...)
}

func (l *logrusWrapper) Fatal(ctx context.Context, args ...interface{}) {
	l.logWithTrace(ctx, logrus.FatalLevel, args...)
	os.Exit(1)
}

func (l *logrusWrapper) Tracer() trace.Tracer {
	return l.tracer
}
