package godotenv

import (
	"github.com/joho/godotenv"
	"os"
)

type Env struct {
	Environment                 string //development,staging,production
	DatabaseHost                string
	DatabaseName                string
	HTTPPort                    string
	Secret                      string
	BaseURL                     string
	RabbitMQHost                string
	RabbitMQEventExchange       string
	RabbitMQInternalExchange    string
	RabbitMQGeodataQueue        string
	MinioHost                   string
	MinioAccessToken            string
	MinioSecret                 string
	MinioProfilePictureBucket   string
	MinioDoctorBucket           string
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
	e.BaseURL = os.Getenv("BASE_URL")

	e.RabbitMQHost = os.Getenv("RABBITMQ_HOST")
	e.RabbitMQEventExchange = os.Getenv("RABBITMQ_EVENT_EXCHANGE")
	e.RabbitMQInternalExchange = os.Getenv("RABBITMQ_INTERNAL_EXCHANGE")

	e.MinioHost = os.Getenv("MINIO_HOST")
	e.MinioAccessToken = os.Getenv("MINIO_ACCESS_TOKEN")
	e.MinioSecret = os.Getenv("MINIO_SECRET")
	e.MinioProfilePictureBucket = os.Getenv("MINIO_PROFILE_PICTURE_BUCKET")
	e.MinioDoctorBucket = os.Getenv("MINIO_DOCTOR_BUCKET")
	e.RabbitMQGeodataQueue = os.Getenv("RABBITMQ_GEODATA_QUEUE")

	e.OpenTelemetryMetricExporter = os.Getenv("OPEN_TELEMETRY_METRIC_EXPORTER")
	e.OpenTelemetryLogExporter = os.Getenv("OPEN_TELEMETRY_LOG_EXPORTER")
}
