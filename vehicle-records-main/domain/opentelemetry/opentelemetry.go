package opentelemetry

import (
	"context"
	"errors"

	"git.abanppc.com/farin-project/vehicle-records/infrastructure/godotenv"
	"go.opentelemetry.io/contrib/instrumentation/host"
	"go.opentelemetry.io/contrib/instrumentation/runtime"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/sdk/log"

	"time"

	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

type shutdownFunc func(context.Context) error

type OpenTelemetry struct {
	shutdownFunc  shutdownFunc
	env           *godotenv.Env
	shutdownFuncs []func(context.Context) error
	Logger        *log.LoggerProvider
	Meter         *metric.MeterProvider
}

func NewOpenTelemetry(env *godotenv.Env) *OpenTelemetry {
	return &OpenTelemetry{
		env: env,
	}
}

func (t *OpenTelemetry) Setup(ctx context.Context) error {
	sf, err := t.setupOTelSDK(ctx)
	if err != nil {
		return err
	}
	t.shutdownFunc = sf
	return nil
}

func (t *OpenTelemetry) Shutdown(ctx context.Context) error {
	return t.shutdownFunc(ctx)
}

// setupOTelSDK bootstraps the OpenTelemetry pipeline.
// If it does not return an error, make sure to call shutdown for proper cleanup.
func (t *OpenTelemetry) setupOTelSDK(ctx context.Context) (shutdown shutdownFunc, err error) {
	// shutdown calls cleanup functions registered via shutdownFuncs.
	// The errors from the calls are joined.
	// Each registered cleanup will be invoked once.
	shutdown = func(ctx context.Context) error {
		var err error
		for _, fn := range t.shutdownFuncs {
			err = errors.Join(err, fn(ctx))
		}
		t.shutdownFuncs = nil
		return err
	}

	// handleErr calls shutdown for cleanup and makes sure that all errors are returned.
	handleErr := func(inErr error) {
		err = errors.Join(inErr, shutdown(ctx))
	}

	r, err := t.newResource(ctx)
	if err != nil {
		handleErr(err)
		return
	}

	// Set up Meter provider.
	meterProvider, err := t.newMeterProvider(r)
	if err != nil {
		handleErr(err)
		return
	}
	t.shutdownFuncs = append(t.shutdownFuncs, meterProvider.Shutdown)
	otel.SetMeterProvider(meterProvider)

	// Set up Logger provider.
	loggerProvider, err := t.newLoggerProvider(r)
	if err != nil {
		handleErr(err)
		return
	}
	t.Logger = loggerProvider
	t.shutdownFuncs = append(t.shutdownFuncs, loggerProvider.Shutdown)
	global.SetLoggerProvider(loggerProvider)
	err = host.Start(host.WithMeterProvider(otel.GetMeterProvider()))
	if err != nil {
		handleErr(err)
	}
	err = runtime.Start(runtime.WithMeterProvider(otel.GetMeterProvider()))
	if err != nil {
		handleErr(err)
	}
	return
}

func (t *OpenTelemetry) newResource(ctx context.Context) (*resource.Resource, error) {
	r, err := resource.New(
		ctx,
		resource.WithAttributes(semconv.ServiceNameKey.String("my-service")),
		resource.WithProcessPID(),
		resource.WithTelemetrySDK(),
		resource.WithHost(),
		resource.WithTelemetrySDK(),
		resource.WithProcess(),
	)
	return r, err
}

func (t *OpenTelemetry) newMeterProvider(r *resource.Resource) (*metric.MeterProvider, error) {
	metricExporter, err := otlpmetricgrpc.New(context.Background(),
		otlpmetricgrpc.WithEndpoint(t.env.OpenTelemetryMetricExporter),
		otlpmetricgrpc.WithInsecure(),
		otlpmetricgrpc.WithHeaders(
			map[string]string{
				"Authorization": "Basic cm9vdEBleGFtcGxlLmNvbTp1a0wzeUtmeXUzM25qRmdu",
				"organization":  "default",
				"stream-name":   "default",
			},
		),
	)
	if err != nil {
		return nil, err
	}

	meterProvider := metric.NewMeterProvider(
		metric.WithResource(r),
		metric.WithReader(metric.NewPeriodicReader(metricExporter,
			metric.WithInterval(10*time.Second))),
	)
	t.Meter = meterProvider

	return meterProvider, nil
}

func (t *OpenTelemetry) newLoggerProvider(r *resource.Resource) (*log.LoggerProvider, error) {
	logExporter, err := otlploggrpc.New(context.Background(),
		otlploggrpc.WithEndpoint(t.env.OpenTelemetryLogExporter),
		otlploggrpc.WithInsecure(),
		otlploggrpc.WithHeaders(
			map[string]string{
				"Authorization": "Basic cm9vdEBleGFtcGxlLmNvbTp1a0wzeUtmeXUzM25qRmdu",
				"organization":  "default",
				"stream-name":   "default",
			},
		),
	)
	if err != nil {
		return nil, err
	}

	loggerProvider := log.NewLoggerProvider(
		log.WithProcessor(log.NewBatchProcessor(logExporter)),
	)
	return loggerProvider, nil
}
