package handlers

import (
	"context"
	"farin/app/rabbit/consumers"
	"farin/domain/service"
	"farin/infrastructure/opentelemetry"
	"git.abanppc.com/farin-project/vehicle-records/domain/entity"
	"go.opentelemetry.io/otel/metric"
	"log/slog"
	"time"
)

type VehicleRecord struct {
	logger         *slog.Logger
	vs             *service.VehicleRecordService
	meter          metric.Meter
	consumeCounter metric.Int64Counter
	successCounter metric.Int64Counter
	failCounter    metric.Int64Counter
}

func NewVehicleRecord(logger *slog.Logger, vs *service.VehicleRecordService, telemetry *opentelemetry.OpenTelemetry) *VehicleRecord {
	meter := telemetry.Meter.Meter("farin,backend.RabbitMQVehicleRecordHandler")
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
	return &VehicleRecord{logger: logger.With("layer", "RabbitEventHandler"), vs: vs,
		consumeCounter: consumeCounter, successCounter: successCounter, failCounter: failCounter}
}

func (a *VehicleRecord) VehicleRecord(ctx context.Context, data []byte) error {
	lg := a.logger.With("method", "VehicleRecord")
	a.consumeCounter.Add(ctx, 1)

	ctx, cancel := context.WithTimeout(ctx, time.Second*60)
	defer cancel()
	vr := &entity.VehicleRecord{}
	if err := vr.UnmarshalJSON(data); err != nil {
		lg.Warn("failed to unmarshal vehicle record", "error", err)
		a.failCounter.Add(ctx, 1)
		return nil
	}
	vr, err := a.vs.CreateRecord(ctx, vr)
	if err != nil {
		lg.Warn("failed to create record", "error", err)
		a.failCounter.Add(ctx, 1)
		return err
	}
	a.successCounter.Add(ctx, 1)
	lg.Info("vehicle record created successfully", "record_id", vr.RecordID)
	return nil
}

func (a *VehicleRecord) RegisterConsumer(c consumers.Consumer) {
	c.RegisterHandler("farin.vehicles.drivers.event", a.VehicleRecord)
}
