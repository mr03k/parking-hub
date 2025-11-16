package event

import (
	"context"
	v1 "git.abanppc.com/farin-project/event/domain/entity/v1"
)

type EventRepoInterface interface {
	CreateEvent(ctx context.Context, event *v1.Event) error
}
