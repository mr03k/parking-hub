package handlers

import (
	"context"
	"git.abanppc.com/farin-project/state/app/rabbit/consumers"
	"git.abanppc.com/farin-project/state/domain/entity"
	"git.abanppc.com/farin-project/state/domain/service"
	"log/slog"
	"time"
)

type State struct {
	logger *slog.Logger
	ss     *service.StateService
}

func NewState(logger *slog.Logger, ss *service.StateService) *State {
	return &State{logger: logger.With("layer", "RabbitEventHandler"), ss: ss}
}

func (a *State) State(ctx context.Context, data []byte) error {
	lg := a.logger.With("method", "State Store")
	state := &entity.State{}
	if err := state.UnmarshalJSON(data); err != nil {
		lg.Warn("failed to unmarshal vehicle record", "error", err)
		return nil
	}
	state, err := a.ss.CreateState(ctx, state)
	if err != nil {
		lg.Warn("failed to create record", "error", err)
		return err
	}
	lg.Info("state created successfully", "state_id", state.ID)
	time.Sleep(1 * time.Second)
	return nil
}

func (a *State) RegisterConsumer(c *consumers.StateConsumer) {
	c.RegisterHandler("farin.vehicles.drivers.state", a.State)
}
