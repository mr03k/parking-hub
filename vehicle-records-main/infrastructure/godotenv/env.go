package godotenv

import (
	"github.com/joho/godotenv"
	"os"
)

type Env struct {
	Environment               string //development,staging,production
	DatabaseHost              string
	DatabaseName              string
	HTTPPort                  string
	Secret                    string
	BaseURL                   string
	RabbitMQHost              string
	RabbitMQEventExchange     string
	RabbitMQInternalExchange  string
	RabbitMQRecordsQueue      string
	RabbitMQRetryQueue        string
	MinioHost                 string
	MinioAccessToken          string
	MinioSecret               string
	MinioVehicleRecordsBucket string

	TehranToken          string
	TehranLoginURL       string
	TehranStoreRecordURL string

	OpenTelemetryMetricExporter string
	OpenTelemetryLogExporter    string
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
	e.TehranToken = os.Getenv("TEHRAN_TOKEN")
	e.TehranLoginURL = os.Getenv("TEHRAN_LOGIN_URL")
	e.TehranStoreRecordURL = os.Getenv("TEHRAN_STORE_RECORD_URL")
	e.BaseURL = os.Getenv("BASE_URL")
	e.RabbitMQHost = os.Getenv("RABBITMQ_HOST")
	e.RabbitMQEventExchange = os.Getenv("RABBITMQ_EVENT_EXCHANGE")
	e.RabbitMQInternalExchange = os.Getenv("RABBITMQ_INTERNAL_EXCHANGE")
	e.RabbitMQRecordsQueue = os.Getenv("RABBITMQ_RECORDS_QUEUE")
	e.RabbitMQRetryQueue = os.Getenv("RABBITMQ_RETRY_QUEUE")
	e.MinioHost = os.Getenv("MINIO_HOST")
	e.MinioAccessToken = os.Getenv("MINIO_ACCESS_TOKEN")
	e.MinioSecret = os.Getenv("MINIO_SECRET")
	e.MinioVehicleRecordsBucket = os.Getenv("MINIO_VEHICLE_RECORDS_BUCKET")
	e.OpenTelemetryMetricExporter = os.Getenv("OPEN_TELEMETRY_METRIC_EXPORTER")
	e.OpenTelemetryLogExporter = os.Getenv("OPEN_TELEMETRY_LOG_EXPORTER")
}
