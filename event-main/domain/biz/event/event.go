package event

import (
	"context"
	v1 "git.abanppc.com/farin-project/event/domain/entity/v1"
	v12 "git.abanppc.com/farin-project/event/domain/repository/v1/event"
	"go.opentelemetry.io/otel"
	"log/slog"
)

type EventBizInterface interface {
	CreateEvent(ctx context.Context, event *v1.Event) error
}

type EventBiz struct {
	repo   v12.EventRepoInterface
	logger *slog.Logger
}

// New Usecase
func NewEventBiz(repo v12.EventRepoInterface, logger *slog.Logger) EventBizInterface {
	return &EventBiz{
		repo:   repo,
		logger: logger.With("layer", "EventBiz"),
	}
}

func (uc *EventBiz) CreateEvent(ctx context.Context, event *v1.Event) error {
	lg := uc.logger.With("method", "CreateEvent")
	ctx, span := otel.Tracer("biz").Start(ctx, "CreateEvent")
	defer span.End()

	if err := uc.repo.CreateEvent(ctx, event); err != nil {
		lg.Error("Error in CreateEvent Repo", "err", err)
		return err
	}
	return nil
}
