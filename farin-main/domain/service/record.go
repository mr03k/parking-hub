package service

import (
	"context"
	"farin/domain/repository"
	"farin/infrastructure/godotenv"
	"farin/infrastructure/opentelemetry"
	vehicleRecordEnt "git.abanppc.com/farin-project/vehicle-records/domain/entity"
	"github.com/mahdimehrabi/uploader"
	"go.opentelemetry.io/otel/metric"
	"log/slog"
	"time"
)

type VehicleRecordService struct {
	logger             *slog.Logger
	recordRepo         *repository.VehicleRecordRepository
	photoRepo          *repository.CitizenVehiclePhotoRepository
	eventRecordRbtRepo *repository.EventRecordRabbitMQ
	fr                 uploader.FileRepository
	env                *godotenv.Env
	durationCounter    metric.Int64Histogram
}

func NewVehicleRecordService(logger *slog.Logger, ringRepo *repository.VehicleRecordRepository, fr uploader.FileRepository,
	env *godotenv.Env, photoRepo *repository.CitizenVehiclePhotoRepository, telemetry *opentelemetry.OpenTelemetry,
	eventRecordRbtRepo *repository.EventRecordRabbitMQ) *VehicleRecordService {
	meter := telemetry.Meter.Meter("farin.backend.VehicleRecordService")

	durationCounter, err := meter.Int64Histogram("vehicle_record.consume.duration",
		metric.WithDescription("duration of processing data in record service"))
	if err != nil {
		panic(err)
	}

	if err != nil {
		panic(err)
	}
	return &VehicleRecordService{
		logger:             logger.With("layer", "VehicleRecordService"),
		recordRepo:         ringRepo,
		fr:                 fr,
		env:                env,
		photoRepo:          photoRepo,
		eventRecordRbtRepo: eventRecordRbtRepo,
		durationCounter:    durationCounter,
	}
}

func (s *VehicleRecordService) CreateRecord(ctx context.Context, record *vehicleRecordEnt.VehicleRecord) (*vehicleRecordEnt.VehicleRecord, error) {
	lg := s.logger.With("method", "CreateRecord")
	startTime := time.Now().UnixMicro()

	rings, err := s.recordRepo.FindRing(ctx, record)
	if err != nil {
		lg.Error("failed to find record", "error", err)
		return nil, err
	}
	if len(rings) < 1 {
		record.RingID = 0
	} else {
		record.RingID = rings[0].ID
	}

	segments, err := s.recordRepo.FindSegment(ctx, record)
	if err != nil {
		lg.Error("failed to find record", "error", err)
		return nil, err
	}
	if len(segments) < 1 {
		record.SegmentID = 0
	} else {
		segment := segments[0]
		record.SegmentID = segment.ID
		record.IsJunction = segment.Junction == 1
	}

	roads, err := s.recordRepo.FindRoad(ctx, record)
	if err != nil {
		lg.Error("failed to find record", "error", err)
		return nil, err
	}
	if len(roads) < 1 {
		record.StreetID = 0
	} else {
		road := roads[0]
		record.StreetID = road.ID
		record.RoadCode = road.RoadCode
	}

	parkings, err := s.recordRepo.FindParking(ctx, record)
	if err != nil {
		lg.Error("failed to find record", "error", err)
		return nil, err
	}
	if len(parkings) < 1 {
		record.ParkingLotID = 0
	} else {
		record.ParkingLotID = parkings[0].ID
	}

	if err := s.eventRecordRbtRepo.Store(ctx, record); err != nil {
		lg.Error("failed to store record", "error", err)
		return nil, err
	}
	endTime := time.Now().UnixMicro()
	s.durationCounter.Record(ctx, endTime-startTime)

	lg.Info("record created", "ringID", record.RecordID)
	return record, nil
}
