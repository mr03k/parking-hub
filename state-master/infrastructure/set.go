package infrastructure

import (
	"git.abanppc.com/farin-project/state/infrastructure/godotenv"
	"git.abanppc.com/farin-project/state/infrastructure/rabbit"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	godotenv.NewEnv,
	rabbit.NewConsumerRunner,
)
