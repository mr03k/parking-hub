package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"git.abanppc.com/farin-project/vehicle-records/domain/entity"
	"git.abanppc.com/farin-project/vehicle-records/infrastructure/godotenv"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// MarginParkingRequest represents the main request structure for margin parking detection
type MarginParkingRequest struct {
	MarginParkingPlateDetection MarginParkingPlateDetectionDTO `json:"marginParkingPalteDetection"`
	VehiclePhotos               []VehiclePhotoDTO              `json:"vehiclePhotos"`
}

// MarginParkingPlateDetectionDTO represents the plate detection data
type MarginParkingPlateDetectionDTO struct {
	RecordID                           string    `json:"recordId"`
	RecordStoreTime                    time.Time `json:"recordStoreTime"`
	RecordSendTime                     time.Time `json:"recordSendTime"`
	CitizenPlateNumber                 string    `json:"citizenPlateNumber"`
	CitizenVehicleType                 int       `json:"citizenVehicleType"`
	CitizenVehiclePlateNumberType      int       `json:"citizenVehiclePlateNumberType"`
	IsCitizenVehicleDistorted          int       `json:"isCitizenVehicleDistorted"`
	IsCitizenVehiclePlateNumberVisible int       `json:"isCitizenVehiclePlateNumberVisible"`
	CitizenVehiclePlateNumberColor     string    `json:"citizenVehiclePlateNumberColor"`
	RingID                             string    `json:"ringId"`
	StreetID                           string    `json:"streetId"`
	SegmentID                          string    `json:"segmentId"`
	ParkinglotID                       string    `json:"parkinglotId"`
	UserID                             string    `json:"userId"`
	LPRVehicleID                       string    `json:"lprvehicleId"`
	LPRSystemID                        string    `json:"lprsystemId"`
	LPRSystemAppID                     string    `json:"lprsystemAppId"`
	LPRSystemAppVersion                string    `json:"lprsystemAppVersion"`
	LPRVehicleGPSLatitude              float64   `json:"lprvehicleGpslatitude"`
	LPRVehicleGPSLongitude             float64   `json:"lprvehicleGpslongitude"`
	LPRVehicleGPSSpeed                 float64   `json:"lprvehicleGpsspeed"`
	LPRVehicleGPSError                 float64   `json:"lprvehicleGpserror"`
	LPRVehicleRTKLatitude              float64   `json:"lprvehicleRtklatitude"`
	LPRVehicleRTKLongitude             float64   `json:"lprvehicleRtklongitude"`
	LPRVehicleRTKError                 float64   `json:"lprvehicleRtkerror"`
	CycleId                            int64     `json:"cycleId"`
	CitizenPlateNumberNumeric          int64     `json:"citizenPlateNumberNumeric"`
	IsJunction                         bool      `json:"isJunction"`
	IsUTC                              bool      `json:"isUtc"`
}

// VehiclePhotoDTO represents the vehicle photo data
type VehiclePhotoDTO struct {
	LPRVehicleCameraID             int       `json:"lprvehicleCameraId"`
	CitizenVehiclePhoto            string    `json:"citizenVehiclePhoto"`
	CitizenVehiclePlateCropPhoto   string    `json:"citizenVehiclePlateCropPhoto"`
	CitizenVehiclePhotoCaptureTime time.Time `json:"citizenVehiclePhotoCaptureTime"`
	OCRAccuracy                    float64   `json:"ocraccuracy"`
}

type AuthTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
}

type MarginParkingResponse struct {
	RequestID        string `json:"requestId"`
	PlateDetectionID int    `json:"plateDetectionId"`
}

// Implementation
type TehranSiteRecordRepository struct {
	httpClient *http.Client
	env        *godotenv.Env
	expire     int64
	token      string
	logger     *slog.Logger
}

// NewTehranSiteRecordRepository creates a new repository instance
func NewTehranSiteRecordRepository(env *godotenv.Env, logger *slog.Logger) *TehranSiteRecordRepository {
	return &TehranSiteRecordRepository{
		httpClient: &http.Client{
			Timeout: 50 * time.Second,
		},
		env:    env,
		logger: logger,
	}
}

// Converter functions
func convertToMarginParkingRequest(record *entity.VehicleRecord) *MarginParkingRequest {
	return &MarginParkingRequest{
		MarginParkingPlateDetection: convertToPlateDetectionDTO(record),
		VehiclePhotos:               convertToVehiclePhotosDTO(record.VehiclePhotos),
	}
}

func convertToPlateDetectionDTO(record *entity.VehicleRecord) MarginParkingPlateDetectionDTO {
	l, err := time.LoadLocation("Asia/Tehran")
	if err != nil {
		l = time.FixedZone("Iran Standard Time", 3*60*60+30*60)
	}

	return MarginParkingPlateDetectionDTO{
		RecordID:                           record.RecordID,
		RecordStoreTime:                    time.Unix(0, record.RecordStoreTime*int64(time.Millisecond)).In(l),
		RecordSendTime:                     time.Unix(0, record.RecordSendTime*int64(time.Millisecond)).In(l),
		CitizenPlateNumber:                 record.CitizenPlateNumber,
		CitizenVehicleType:                 record.CitizenVehicleType,
		CitizenVehiclePlateNumberType:      record.CitizenVehiclePlateNumberType,
		IsCitizenVehicleDistorted:          boolToInt(record.IsCitizenVehicleDistorted),
		IsCitizenVehiclePlateNumberVisible: boolToInt(record.IsCitizenVehiclePlateNumberVisible),
		CitizenVehiclePlateNumberColor:     fmt.Sprint(record.CitizenVehiclePlateNumberColor),
		RingID:                             strconv.FormatInt(record.RingID, 10),
		StreetID:                           strconv.FormatInt(record.RoadCode, 10),
		SegmentID:                          strconv.FormatInt(record.SegmentID, 10),
		ParkinglotID:                       strconv.FormatInt(record.ParkingLotID, 10),
		UserID:                             record.UserID,
		LPRVehicleID:                       record.LPRVehicleID,
		LPRSystemID:                        record.LPRSystemID,
		LPRSystemAppID:                     record.LPRSystemAppID,
		LPRSystemAppVersion:                record.LPRSystemAppVersion,
		LPRVehicleGPSLatitude:              record.LPRVehicleGPSLatitude,
		LPRVehicleGPSLongitude:             record.LPRVehicleGPSLongitude,
		LPRVehicleGPSSpeed:                 record.LPRVehicleGPSSpeed,
		LPRVehicleGPSError:                 float64(record.LPRVehicleGPSError),
		LPRVehicleRTKLatitude:              record.LPRVehicleRTKLatitude,
		LPRVehicleRTKLongitude:             record.LPRVehicleRTKLongitude,
		LPRVehicleRTKError:                 float64(record.LPRVehicleRTKError),
		CycleId:                            int64(record.CycleID),
		CitizenPlateNumberNumeric:          int64(record.CitizenPlateNumberNumeric),
		IsJunction:                         record.IsJunction,
		IsUTC:                              false,
	}
}

func convertToVehiclePhotosDTO(photos []*entity.CitizenVehiclePhoto) []VehiclePhotoDTO {
	l, err := time.LoadLocation("Asia/Tehran")
	if err != nil {
		l = time.FixedZone("Iran Standard Time", 3*60*60+30*60)
	}
	result := make([]VehiclePhotoDTO, len(photos))
	for i, photo := range photos {
		result[i] = VehiclePhotoDTO{
			LPRVehicleCameraID:           photo.LPRVehicleCameraID,
			CitizenVehiclePhoto:          photo.CitizenVehiclePhoto,
			CitizenVehiclePlateCropPhoto: photo.CitizenVehiclePlateCropPhoto,
			CitizenVehiclePhotoCaptureTime: time.Unix(0,
				photo.CitizenVehiclePhotoCaptureTime*int64(time.Millisecond)).In(l),
			OCRAccuracy: photo.OCRAccuracy,
		}
	}
	return result
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
func (r *TehranSiteRecordRepository) Authenticate(ctx context.Context) error {
	// Create form data
	formData := url.Values{}
	formData.Set("grant_type", "client_credentials")
	formData.Set("Scope", "")

	// Create request
	req, err := http.NewRequestWithContext(ctx,
		http.MethodPost,
		r.env.TehranLoginURL,
		strings.NewReader(formData.Encode()))
	if err != nil {
		return fmt.Errorf("failed to create auth request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Basic "+r.env.TehranToken)

	// Send request
	resp, err := r.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send auth request: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("authentication failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	// Parse response
	var tokenResponse AuthTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		return fmt.Errorf("failed to decode auth response: %w", err)
	}
	r.expire = time.Now().Add(time.Second * time.Duration(tokenResponse.ExpiresIn)).Unix()
	r.token = tokenResponse.AccessToken
	return nil
}

// SendVehicleRecord sends the record to the external service
func (r *TehranSiteRecordRepository) SendVehicleRecord(ctx context.Context, record *entity.VehicleRecord) (requestID string, plateDetctionID int, err error) {
	lg := r.logger.With("method", "SendVehicleRecord")
	if time.Now().Unix() > r.expire-10 {
		if err = r.Authenticate(ctx); err != nil {
			return
		}
	}
	requestDTO := convertToMarginParkingRequest(record)

	jsonData, err := json.Marshal(requestDTO)
	if err != nil {
		err = fmt.Errorf("failed to marshal request: %w", err)
		return
	}

	req, err := http.NewRequestWithContext(ctx,
		http.MethodPost,
		r.env.TehranStoreRecordURL,
		bytes.NewBuffer(jsonData))
	if err != nil {
		err = fmt.Errorf("failed to create request: %w", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+r.token)

	//every 50 request log one shardari request
	if time.Now().Unix()%50 == 0 {
		lg.Info("shardari request", slog.String("req", string(jsonData)))
	}
	resp, err := r.httpClient.Do(req)
	if err != nil {
		err = fmt.Errorf("failed to send request: %w", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var b []byte
		b, err = io.ReadAll(resp.Body)
		if err != nil {
			err = fmt.Errorf("failed to read response body: %w", err)
			return
		}
		err = fmt.Errorf("unexpected status code: %d and resp:%s", resp.StatusCode, string(b), "Req", string(jsonData))
		return
	}

	// Parse response
	var response MarginParkingResponse
	if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
		err = fmt.Errorf("failed to decode response: %w", err)
		return
	}
	requestID = response.RequestID
	plateDetctionID = response.PlateDetectionID

	return
}
