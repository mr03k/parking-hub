package consumers

import (
	"git.abanppc.com/farin-project/event/config"
	"git.abanppc.com/farin-project/event/domain/biz/event"
	"git.abanppc.com/farin-project/event/pkg/rabbit"
	"log/slog"
)

type Consumer interface {
	Setup() error
	Consume()
}

func Consumers(logger *slog.Logger, rbt *rabbit.Rabbit, rbtConfig *config.RabbitMQ, eventBiz event.EventBizInterface) []Consumer {
	return []Consumer{
		NewEvent(logger, rbt, rbtConfig, eventBiz),
	}
}
