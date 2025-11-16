package infrastructure

import (
	"git.abanppc.com/farin-project/vehicle-records/infrastructure/godotenv"
	"git.abanppc.com/farin-project/vehicle-records/infrastructure/rabbit"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	godotenv.NewEnv,
	rabbit.NewConsumerRunner,
)
