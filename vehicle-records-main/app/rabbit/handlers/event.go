package handlers

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"git.abanppc.com/farin-project/vehicle-records/app/rabbit/consumers"
	"git.abanppc.com/farin-project/vehicle-records/domain/entity"
	"git.abanppc.com/farin-project/vehicle-records/domain/opentelemetry"
	"git.abanppc.com/farin-project/vehicle-records/domain/repository"
	"git.abanppc.com/farin-project/vehicle-records/domain/service"
	"go.opentelemetry.io/otel/metric"
)

type EventVehicleRecord struct {
	logger         *slog.Logger
	vs             *service.VehicleRecordService
	consumeCounter metric.Int64Counter
	successCounter metric.Int64Counter
	failCounter    metric.Int64Counter
}

func NewEventVehicleRecord(logger *slog.Logger, vs *service.VehicleRecordService, telemetry *opentelemetry.OpenTelemetry) *EventVehicleRecord {
	meter := telemetry.Meter.Meter("farin,vehicleRecord.RabbitMQVehicleRecordHandler")
	consumeCounter, err := meter.Int64Counter("vehicle_record.consume.total.counter", metric.WithDescription("number of vehicle record handler consumes"))
	if err != nil {
		panic(err)
	}
	successCounter, err := meter.Int64Counter("vehicle_record.consume.counter", metric.WithDescription("number of vehicle record handler success consumes"))
	if err != nil {
		panic(err)
	}
	failCounter, err := meter.Int64Counter("vehicle_record.consume.failed.counter", metric.WithDescription("number of vehicle record handler failed consumes"))
	if err != nil {
		panic(err)
	}
	return &EventVehicleRecord{logger: logger.With("layer", "RabbitEventHandler"), vs: vs,
		consumeCounter: consumeCounter, successCounter: successCounter, failCounter: failCounter}
}

func (a *EventVehicleRecord) VehicleRecord(ctx context.Context, data []byte) error {
	lg := a.logger.With("method", "VehicleRecord")
	dt := time.Now()
	a.consumeCounter.Add(ctx, 1)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*100)
	defer func() {
		lg.Info("vehicle record handler finished", "duration", time.Since(dt).Seconds())
		cancel()
	}()

	vr := &entity.VehicleRecord{}
	if err := vr.UnmarshalJSON(data); err != nil {
		lg.Warn("failed to unmarshal vehicle record", "error", err)
		a.failCounter.Add(ctx, 1)
		return nil
	}
	cr, err := a.vs.CreateRecord(ctx, vr)
	if err != nil {
		if errors.Is(err, repository.ErrDuplicateRecord) {
			lg.Warn("already exists vehicle record", "error", err, "record_id", vr.RecordID)
			a.failCounter.Add(ctx, 1)
			return nil
		}
		lg.Warn("failed to create record", "error", err)
		a.failCounter.Add(ctx, 1)
		return err
	}
	lg.Info("vehicle record created successfully", "record_id", cr.RecordID)
	a.successCounter.Add(ctx, 1)
	return nil
}

func (a *EventVehicleRecord) RegisterConsumer(c *consumers.EventConsumer) {
	c.RegisterHandler("farin.vehicles.drivers.event", a.VehicleRecord)
}
