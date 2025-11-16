package handlers

import (
	"context"
	"log/slog"
	"time"

	"git.abanppc.com/farin-project/vehicle-records/app/rabbit/consumers"
	"git.abanppc.com/farin-project/vehicle-records/domain/entity"
	"git.abanppc.com/farin-project/vehicle-records/domain/opentelemetry"
	"git.abanppc.com/farin-project/vehicle-records/domain/service"
	"go.opentelemetry.io/otel/metric"
)

type Retry struct {
	logger         *slog.Logger
	vs             *service.VehicleRecordService
	consumeCounter metric.Int64Counter
	successCounter metric.Int64Counter
	failCounter    metric.Int64Counter
}

func NewRetry(logger *slog.Logger, vs *service.VehicleRecordService, telemetry *opentelemetry.OpenTelemetry) *Retry {
	meter := telemetry.Meter.Meter("farin,vehicle_record.RabbitMQVehicleRecordHandler")
	consumeCounter, err := meter.Int64Counter("vehicle_record.consume.total.counter",
		metric.WithDescription("number of vehicle record retry handler consumes"))
	if err != nil {
		panic(err)
	}
	successCounter, err := meter.Int64Counter("vehicle_record.retry.consume.counter",
		metric.WithDescription("number of vehicle record retry handler success consumes"))
	if err != nil {
		panic(err)
	}
	failCounter, err := meter.Int64Counter("vehicle_record.retry.consume.failed.counter",
		metric.WithDescription("number of vehicle record  retry handler failed consumes"))
	if err != nil {
		panic(err)
	}

	return &Retry{logger: logger.With("layer", "RabbitEventHandler"), vs: vs,
		consumeCounter: consumeCounter, successCounter: successCounter, failCounter: failCounter}
}

func (a *Retry) Retry(ctx context.Context, data []byte) error {
	lg := a.logger.With("method", "Retry")

	dt := time.Now()
	defer func() {
		lg.Info("vehicle record retry handler finished", "duration", time.Since(dt).Seconds())
	}()

	a.consumeCounter.Add(ctx, 1)
	ctx, cancell := context.WithTimeout(context.Background(), time.Second*100)
	defer func() {
		lg.Info("vehicle record retry handler finished", "duration", time.Since(dt).Seconds())
		cancell()
	}()
	vr := &entity.VehicleRecord{}
	if err := vr.UnmarshalJSON(data); err != nil {
		lg.Warn("failed to unmarshal vehicle record", "error", err)
		a.failCounter.Add(ctx, 1)
		return nil
	}
	if err := a.vs.Retry(ctx, vr); err != nil {
		lg.Warn("failed to retry create record", "error", err)
		a.failCounter.Add(ctx, 1)
		return err
	}
	a.successCounter.Add(ctx, 1)
	lg.Info("vehicle record created successfully", "record_id", vr.RecordID)
	return nil
}

func (a *Retry) RegisterConsumer(c *consumers.RetryConsumer) {
	c.RegisterHandler("farin.vehicles.drivers.event.retry", a.Retry)
}
