package infrastructure

import (
	"farin/infrastructure/godotenv"
	"farin/infrastructure/rabbit"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	godotenv.NewEnv,
	rabbit.NewConsumerRunner,
)
