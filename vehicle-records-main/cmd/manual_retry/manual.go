package manual_retry

import (
	"context"
	"git.abanppc.com/farin-project/vehicle-records/domain/service"
	"git.abanppc.com/farin-project/vehicle-records/infrastructure/godotenv"
	"log"
	"log/slog"
)

type ManualRetry struct {
	log           *slog.Logger
	env           *godotenv.Env
	recordService *service.VehicleRecordService
}

func NewManualRetry(log *slog.Logger, env *godotenv.Env, recordService *service.VehicleRecordService) *ManualRetry {
	return &ManualRetry{log: log, env: env, recordService: recordService}
}

func (r ManualRetry) FindResend(ctx context.Context, fromTime int64, limit int) {
	if err := r.recordService.FindResend(ctx, fromTime, limit); err != nil {
		log.Fatal(err.Error())
	}
}
