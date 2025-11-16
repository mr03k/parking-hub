package godotenv

import (
	"github.com/joho/godotenv"
	"os"
)

type Env struct {
	Environment           string //development,staging,production
	DatabaseHost          string
	DatabaseName          string
	HTTPPort              string
	Secret                string
	BaseURL               string
	RabbitMQHost          string
	RabbitMQEventExchange string
	RabbitMQStateQueue    string
}

func NewEnv() *Env {
	e := &Env{}
	e.Load()
	return e
}

func (e *Env) Load() {
	godotenv.Load(".env") // using .env file is not mandatory
	e.Environment = os.Getenv("ENVIRONMENT")
	e.DatabaseHost = os.Getenv("DATABASE_HOST")
	e.DatabaseName = os.Getenv("DATABASE_NAME")
	e.HTTPPort = os.Getenv("HTTP_PORT")
	e.Secret = os.Getenv("SECRET")
	e.BaseURL = os.Getenv("BASE_URL")
	e.RabbitMQHost = os.Getenv("RABBITMQ_HOST")
	e.RabbitMQEventExchange = os.Getenv("RABBITMQ_EVENT_EXCHANGE")
	e.RabbitMQStateQueue = os.Getenv("RABBITMQ_STATE_QUEUE")
}
