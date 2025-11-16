package infrastructure

import (
	"application/config"
	"context"
	"errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"time"
)

type OpenTelemetry struct {
	obsConfig config.Observability
}

func NewOpenTelemetry(obsConfig config.Observability) *OpenTelemetry {
	return &OpenTelemetry{
		obsConfig: obsConfig,
	}
}

func (t OpenTelemetry) Setup() error {
	ctx := context.Background()

	// Setting up the OpenTelemetry components
	shutdown, err := t.setupOTelSDK(ctx)
	if err != nil {
		return err
	}
	defer shutdown(ctx)

	return nil
}

// setupOTelSDK bootstraps the OpenTelemetry pipeline.
func (t OpenTelemetry) setupOTelSDK(ctx context.Context) (shutdown func(context.Context) error, err error) {
	var shutdownFuncs []func(context.Context) error

	// shutdown calls cleanup functions registered via shutdownFuncs.
	shutdown = func(ctx context.Context) error {
		var err error
		for _, fn := range shutdownFuncs {
			err = errors.Join(err, fn(ctx))
		}
		shutdownFuncs = nil
		return err
	}

	// handleErr calls shutdown for cleanup and makes sure that all errors are returned.
	handleErr := func(inErr error) {
		err = errors.Join(inErr, shutdown(ctx))
	}

	// Set up propagator.
	prop := t.newPropagator()
	otel.SetTextMapPropagator(prop)

	// Set up trace provider.
	tracerProvider, err := t.newTraceProvider()
	if err != nil {
		handleErr(err)
		return
	}
	shutdownFuncs = append(shutdownFuncs, tracerProvider.Shutdown)
	otel.SetTracerProvider(tracerProvider)

	// Set up meter provider.
	meterProvider, err := t.newMeterProvider()
	if err != nil {
		handleErr(err)
		return
	}
	shutdownFuncs = append(shutdownFuncs, meterProvider.Shutdown)
	otel.SetMeterProvider(meterProvider)

	// Set up logger provider.
	loggerProvider, err := t.newLoggerProvider()
	if err != nil {
		handleErr(err)
		return
	}
	shutdownFuncs = append(shutdownFuncs, loggerProvider.Shutdown)
	global.SetLoggerProvider(loggerProvider)
	return
}

func (t OpenTelemetry) newPropagator() propagation.TextMapPropagator {
	return propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
}

func (t OpenTelemetry) newTraceProvider() (*trace.TracerProvider, error) {
	traceExporter, err := otlptracegrpc.New(context.Background(),
		otlptracegrpc.WithEndpoint("localhost:5081"),
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithHeaders(
			map[string]string{
				"Authorization": t.obsConfig.Tracing.AuthorizationHeader,
				"organization":  t.obsConfig.Tracing.Organization,
				"stream-name":   t.obsConfig.Tracing.StreamName,
			}),
	)

	if err != nil {
		return nil, err
	}

	traceProvider := trace.NewTracerProvider(
		trace.WithBatcher(traceExporter,
			trace.WithBatchTimeout(time.Second)), // Default is 5s. Set to 1s for demonstrative purposes.
	)
	return traceProvider, nil
}

func (t OpenTelemetry) newMeterProvider() (*metric.MeterProvider, error) {
	metricExporter, err := otlpmetricgrpc.New(context.Background(),
		otlpmetricgrpc.WithEndpoint("localhost:5081"),
		otlpmetricgrpc.WithInsecure(),
		otlpmetricgrpc.WithHeaders(map[string]string{
			"Authorization": t.obsConfig.Metrics.AuthorizationHeader,
			"organization":  t.obsConfig.Metrics.Organization,
			"stream-name":   t.obsConfig.Metrics.StreamName,
		}),
	)
	if err != nil {
		return nil, err
	}
	r, err := resource.New(
		context.Background(),
		resource.WithAttributes(semconv.ServiceNameKey.String("my-service")),
		resource.WithProcessPID(),
		resource.WithTelemetrySDK(),
		resource.WithHost(),
		resource.WithTelemetrySDK(),
		resource.WithProcess(),
	)

	meterProvider := metric.NewMeterProvider(
		metric.WithResource(r),
		metric.WithReader(metric.NewPeriodicReader(metricExporter,
			metric.WithInterval(3*time.Second))), // Default is 1m. Set to 3s for demonstrative purposes.
	)
	return meterProvider, nil
}

func (t OpenTelemetry) newLoggerProvider() (*log.LoggerProvider, error) {
	logExporter, err := otlploggrpc.New(context.Background(),
		otlploggrpc.WithEndpoint("localhost:5081"),
		otlploggrpc.WithHeaders(
			map[string]string{
				"Authorization": t.obsConfig.Logging.AuthorizationHeader,
				"organization":  t.obsConfig.Logging.Organization,
				"stream-name":   t.obsConfig.Logging.StreamName,
			},
		),
		otlploggrpc.WithInsecure(),
	)
	if err != nil {
		return nil, err
	}

	loggerProvider := log.NewLoggerProvider(
		log.WithProcessor(log.NewBatchProcessor(logExporter)),
	)
	return loggerProvider, nil
}
