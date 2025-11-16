package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/rabbitmq/amqp091-go"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"
)

// Event struct
type Event struct {
	ID                                 string  `json:"id" bson:"_id"`
	RecordID                           string  `bson:"recordId,omitempty" json:"recordId"`
	RecordStoreTime                    int64   `bson:"recordStoreTime,omitempty" json:"recordStoreTime"`
	RecordSendTime                     int64   `bson:"recordSendTime,omitempty" json:"recordSendTime"`
	CitizenVehicleType                 int     `bson:"citizenVehicleType,omitempty" json:"citizenVehicleType"`
	CitizenVehicleColor                int     `bson:"citizenVehicleColor,omitempty" json:"citizenVehicleColor"`
	CitizenVehicleModel                string  `bson:"citizenVehicleModel,omitempty" json:"citizenVehicleModel"`
	CitizenVehicleDistance             int     `bson:"citizenVehicleDistance,omitempty" json:"citizenVehicleDistance"`
	CitizenVehicleDegree               int     `bson:"citizenVehicleDegree,omitempty" json:"citizenVehicleDegree"`
	CitizenPlateNumber                 string  `bson:"citizenPlateNumber,omitempty" json:"citizenPlateNumber"`
	CitizenVehiclePlateNumberType      int     `bson:"citizenVehiclePlateNumberType,omitempty" json:"citizenVehiclePlateNumberType"`
	CitizenVehiclePlateNumberColor     int     `bson:"citizenVehiclePlateNumberColor,omitempty" json:"citizenVehiclePlateNumberColor"`
	OCRAccuracy                        float64 `bson:"ocrAccuracy,omitempty" json:"ocrAccuracy"`
	IsCitizenVehicleDistorted          bool    `bson:"isCitizenVehicleDistorted,omitempty" json:"isCitizenVehicleDistorted"`
	IsCitizenVehiclePlateNumberVisible bool    `bson:"isCitizenVehiclePlateNumberVisible,omitempty" json:"isCitizenVehiclePlateNumberVisible"`
	CitizenParkType                    int     `bson:"citizenParkType,omitempty" json:"citizenParkType"`
	RingID                             string  `bson:"ringId,omitempty" json:"ringId"`
	StreetID                           string  `bson:"streetId,omitempty" json:"streetId"`
	UserID                             string  `bson:"userId,omitempty" json:"userId"`
	LPRVehicleID                       string  `bson:"lprVehicleId,omitempty" json:"lprVehicleId"`
	LPRSystemID                        string  `bson:"lprSystemId,omitempty" json:"lprSystemId"`
	LPRSystemAppID                     string  `bson:"lprSystemAppId,omitempty" json:"lprSystemAppId"`
	LPRSystemAppVersion                string  `bson:"lprSystemAppVersion,omitempty" json:"lprSystemAppVersion"`
	LPRSChannelID                      string  `bson:"lprsChaninId,omitempty" json:"lprsChaninId"`
	LPRVehicleGPSSpeed                 float64 `bson:"lprVehicleGPSSpeed,omitempty" json:"lprVehicleGPSSpeed"`
	LPRVehicleIsGPSSignalValid         bool    `bson:"lprVehicleIsGPSSignalValid,omitempty" json:"lprVehicleIsGPSSignalValid"`
	LPRVehicleGPSLocation              Geo     `bson:"lprVehicleGPSLocation,omitempty" json:"lprVehicleGPSLocation"`
	LPRVehicleGPSError                 float64 `bson:"lprVehicleGPSError,omitempty" json:"lprVehicleGPSError"`
	LPRVehicleRTKLocation              Geo     `bson:"lprVehicleRTKLocation,omitempty" json:"lprVehicleRTKLocation"`
	LPRVehicleRTKError                 float64 `bson:"lprVehicleRTKError,omitempty" json:"lprVehicleRTKError"`
	LPRVehicleCameraID                 int     `bson:"lprVehicleCameraId,omitempty" json:"lprVehicleCameraId"`
	CitizenVehiclePhoto                string  `bson:"citizenVehiclePhoto,omitempty" json:"citizenVehiclePhoto"`
	CitizenVehiclePhotoArea            string  `bson:"citizenVehiclePhotoArea,omitempty" json:"citizenVehiclePhotoArea"`
	CitizenVehiclePlateNumberPhoto     string  `bson:"citizenVehiclePlateNumberPhoto,omitempty" json:"citizenVehiclePlateNumberPhoto"`
	CitizenVehiclePhotoCaptureTime     int64   `bson:"citizenVehiclePhotoCaptureTime,omitempty" json:"citizenVehiclePhotoCaptureTime"`
	Published                          string  `bson:"published,omitempty" json:"published"`
}

type Geo struct {
	Latitude  float64 `json:"latitude" bson:"latitude"`
	Longitude float64 `json:"longitude" bson:"longitude"`
}

// Helper function to handle errors
func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	// Open the CSV file
	file, err := os.Open("atoor.csv")
	failOnError(err, "Failed to open file")
	defer file.Close()

	// Create a new CSV reader
	reader := csv.NewReader(file)

	// Read the header row (assuming first row is header)
	_, err = reader.Read() // Skip the header
	failOnError(err, "Failed to read headers")

	// Read the remaining rows
	var records [][]string
	for {
		row, err := reader.Read()
		if err != nil {
			break // EOF will trigger here
		}
		records = append(records, row)
	}

	// Set up RabbitMQ connection and channel
	conn, err := amqp091.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"farin",
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to declare an exchange")

	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)

	// Simulate random record publishing
	for {
		time.Sleep(time.Duration(r.Intn(10)) * time.Second)
		row := records[r.Intn(len(records))]

		// Create an Event struct and populate fields from CSV row
		event := Event{
			RecordID:                           row[0],
			RecordStoreTime:                    parseInt64(row[1]),
			RecordSendTime:                     parseInt64(row[2]),
			CitizenVehicleType:                 parseInt(row[3]),
			CitizenVehicleColor:                parseInt(row[4]),
			CitizenVehicleModel:                row[5],
			CitizenVehicleDistance:             parseInt(row[6]),
			CitizenVehicleDegree:               parseInt(row[7]),
			CitizenPlateNumber:                 row[8],
			CitizenVehiclePlateNumberType:      parseInt(row[9]),
			CitizenVehiclePlateNumberColor:     parseInt(row[10]),
			OCRAccuracy:                        parseFloat(row[11]),
			IsCitizenVehicleDistorted:          parseBool(row[12]),
			IsCitizenVehiclePlateNumberVisible: parseBool(row[13]),
			CitizenParkType:                    parseInt(row[14]),
			RingID:                             row[15],
			StreetID:                           row[16],
			UserID:                             row[17],
			LPRVehicleID:                       row[18],
			LPRSystemID:                        row[19],
			LPRSystemAppID:                     row[20],
			LPRSystemAppVersion:                row[21],
			LPRSChannelID:                      row[22],
			LPRVehicleGPSSpeed:                 parseFloat(row[23]),
			LPRVehicleIsGPSSignalValid:         parseBool(row[24]),
			LPRVehicleGPSLocation: Geo{
				Latitude:  parseFloat(row[25]),
				Longitude: parseFloat(row[26]),
			},
			LPRVehicleGPSError:             parseFloat(row[27]),
			LPRVehicleRTKLocation:          Geo{Latitude: parseFloat(row[28]), Longitude: parseFloat(row[29])},
			LPRVehicleRTKError:             parseFloat(row[30]),
			LPRVehicleCameraID:             parseInt(row[31]),
			CitizenVehiclePhoto:            row[32],
			CitizenVehiclePhotoArea:        row[33],
			CitizenVehiclePlateNumberPhoto: row[34],
			CitizenVehiclePhotoCaptureTime: parseInt64(row[35]),
			Published:                      time.Now().Format(time.RFC3339),
		}

		// Marshal the Event struct to JSON
		js, err := json.Marshal(event)
		failOnError(err, "Failed to marshal event")

		// Publish to RabbitMQ
		err = ch.Publish(
			"farin",
			fmt.Sprintf("farin.vehicles.%s.drivers.%s.event", event.LPRVehicleID, event.UserID),
			false,
			false,
			amqp091.Publishing{
				ContentType: "application/json",
				Body:        js,
			},
		)
		failOnError(err, "Failed to publish a message")

		fmt.Println("Published event for plate number:", event.CitizenPlateNumber)
	}
}

// parseInt converts a string to an int
func parseInt(value string) int {
	i, err := strconv.Atoi(value)
	if err != nil {
		log.Printf("Error parsing int from value: %s, defaulting to 0", value)
		return 0
	}
	return i
}

// parseInt64 converts a string to an int64
func parseInt64(value string) int64 {
	i, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		log.Printf("Error parsing int64 from value: %s, defaulting to 0", value)
		return 0
	}
	return i
}

// parseFloat converts a string to a float64
func parseFloat(value string) float64 {
	f, err := strconv.ParseFloat(value, 64)
	if err != nil {
		log.Printf("Error parsing float64 from value: %s, defaulting to 0.0", value)
		return 0.0
	}
	return f
}

// parseBool converts a string to a boolean
func parseBool(value string) bool {
	b, err := strconv.ParseBool(value)
	if err != nil {
		log.Printf("Error parsing bool from value: %s, defaulting to false", value)
		return false
	}
	return b
}
