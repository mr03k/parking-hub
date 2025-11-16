package entity

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/valyala/fastjson"
	"gorm.io/plugin/soft_delete"
)

type VehicleRecord struct {
	RecordID string `gorm:"type:varchar;primaryKey" json:"RecordId"`

	RecordStoreTime int64 `gorm:"column:record_store_time;not null" json:"RecordStoreTime"`
	RecordSendTime  int64 `gorm:"column:record_send_time;not null" json:"RecordSendTime"`

	CitizenVehicleType             int    `gorm:"column:citizen_vehicle_type;not null" json:"CitizenVehicleType"`
	CitizenVehicleColor            int    `gorm:"column:citizen_vehicle_color;not null" json:"CitizenVehicleColor"`
	CitizenVehicleModel            string `gorm:"column:citizen_vehicle_model" json:"CitizenVehicleModel"`
	CitizenVehicleDistance         int    `gorm:"column:citizen_vehicle_distance;not null" json:"CitizenVehicleDistance"`
	CitizenVehicleDegree           int    `gorm:"column:citizen_vehicle_degree;not null" json:"CitizenVehicleDegree"`
	CitizenPlateNumber             string `gorm:"column:citizen_plate_number;not null" json:"CitizenPlateNumber"`
	CitizenVehiclePlateNumberType  int    `gorm:"column:citizen_vehicle_plate_number_type;not null" json:"CitizenVehiclePlateNumberType"`
	CitizenVehiclePlateNumberColor int    `gorm:"column:citizen_vehicle_plate_number_color;not null" json:"CitizenVehiclePlateNumberColor"`

	OCRAccuracy                        float64 `gorm:"column:ocr_accuracy;type:numeric(6,6);not null" json:"OCRAccuracy"`
	IsCitizenVehicleDistorted          bool    `gorm:"column:is_citizen_vehicle_distorted;not null" json:"IsCitizenVehicleDistorted"`
	IsCitizenVehiclePlateNumberVisible bool    `gorm:"column:is_citizen_vehicle_plate_number_visible;not null" json:"IsCitizenVehiclePlateNumberVisible"`
	Sent                               bool    `gorm:"column:sent;not null" json:"Sent"`
	Retries                            int     `gorm:"column:retries;not null" json:"Retries"`
	BackoffTime                        int64   `gorm:"-;not null" json:"BackoffTime"` //unix timestamp backoff time
	CitizenParkType                    int     `gorm:"column:citizen_park_type;not null" json:"CitizenParkType"`

	RingID       int64 `gorm:"column:ring_id" json:"RingId"`
	StreetID     int64 `gorm:"column:street_id" json:"StreetId"`
	SegmentID    int64 `gorm:"column:segment_id" json:"SegmentId"`
	ParkingLotID int64 `gorm:"column:parking_lot_id" json:"ParkingLotId"`
	RoadCode     int64 `gorm:"column:road_code" json:"RoadCode"` //streetCode
	IsJunction   bool  `gorm:"column:is_junction;not null" json:"IsJunction"`

	UserID              string `gorm:"type:uuid;column:user_id;not null" json:"UserId"`
	LPRVehicleID        string `gorm:"type:uuid;column:lpr_vehicle_id;not null" json:"LPRVehicleId"`
	LPRSystemID         string `gorm:"column:lpr_system_id;not null" json:"LPRSystemId"`
	LPRSystemAppID      string `gorm:"type:uuid;column:lpr_system_app_id;not null" json:"LPRSystemAppId"`
	LPRSystemAppVersion string `gorm:"column:lpr_system_app_version;not null" json:"LPRSystemAppVersion"`

	LPRVehicleGPSSpeed         float64 `gorm:"column:lpr_vehicle_gps_speed;type:numeric(10,2);default:0" json:"LPRVehicleGPSSpeed"`
	LPRVehicleIsGPSSignalValid bool    `gorm:"column:lpr_vehicle_is_gps_signal_valid;not null" json:"LPRVehicleIsGPSSignalValid"`
	LPRVehicleGPSLatitude      float64 `gorm:"column:lpr_vehicle_gps_latitude;type:numeric(10,8);not null" json:"LPRVehicleGPSLatitude"`
	LPRVehicleGPSLongitude     float64 `gorm:"column:lpr_vehicle_gps_longitude;type:numeric(11,8);not null" json:"LPRVehicleGPSLongitude"`
	LPRVehicleGPSError         int     `gorm:"column:lpr_vehicle_gps_error;not null" json:"LPRVehicleGPSError"`

	LPRVehicleRTKLatitude     float64 `gorm:"column:lpr_vehicle_rtk_latitude;type:numeric(10,8);not null" json:"LPRVehicleRTKLatitude"`
	LPRVehicleRTKLongitude    float64 `gorm:"column:lpr_vehicle_rtk_longitude;type:numeric(11,8);not null" json:"LPRVehicleRTKLongitude"`
	LPRVehicleRTKError        int     `gorm:"column:lpr_vehicle_rtk_error;not null" json:"LPRVehicleRTKError"`
	TehranRequestID           string  `json:"TehranRequestID"`
	PlateDetectionID          int     `json:"PlateDetectionID"`
	CycleID                   int     `json:"CycleID"`
	CitizenPlateNumberNumeric int     `json:"CitizenPlateNumberNumeric"`
	ShamsiTime                string  `json:"ShamsiTime"`

	VehiclePhotos []*CitizenVehiclePhoto `gorm:"foreignKey:RecordID" json:"VehiclePhotos,omitempty"`
	CreatedAt     int64                  `gorm:"autoCreateTime:milli" json:"CreatedAt"`
	UpdatedAt     int64                  `gorm:"autoUpdateTime:milli" json:"UpdatedAt"`
	DeletedAt     soft_delete.DeletedAt  `gorm:"index" json:"-"`
}

type CitizenVehiclePhoto struct {
	Base
	RecordID                       string  `gorm:"type:uuid;column:record_id;not null" json:"RecordId"`
	PhotoSequenceID                int     `gorm:"column:photo_sequence_id;not null" json:"PhotoSequenceId"`
	OCRAccuracy                    float64 `gorm:"column:ocr_accuracy;type:numeric(6,6);not null" json:"OCRAccuracy"`
	LPRVehicleCameraID             int     `gorm:"column:lpr_vehicle_camera_id;not null" json:"LPRVehicleCameraId"`
	CitizenVehiclePhoto            string  `gorm:"column:citizen_vehicle_photo" json:"CitizenVehiclePhoto"`
	CitizenVehiclePhotoArea        string  `gorm:"column:citizen_vehicle_photo_area" json:"CitizenVehiclePhotoArea"`
	CitizenVehiclePlateCropPhoto   string  `gorm:"column:citizen_vehicle_plate_crop_photo" json:"CitizenVehiclePlateCropPhoto"`
	CitizenVehiclePhotoCaptureTime int64   `gorm:"column:citizen_vehicle_photo_capture_time;not null" json:"CitizenVehiclePhotoCaptureTime"`

	VehicleRecord VehicleRecord `gorm:"foreignKey:RecordID" json:"-"`
}

func (vr *VehicleRecord) UnmarshalJSON(data []byte) error {
	var p fastjson.Parser
	v, err := p.ParseBytes(data)
	if err != nil {
		return err
	}

	if recordIDBytes := v.GetStringBytes("RecordId"); len(recordIDBytes) > 0 {
		recordID, err := uuid.Parse(string(recordIDBytes))
		if err != nil {
			return fmt.Errorf("invalid RecordId: %v", err)
		}
		vr.RecordID = recordID.String()
	}

	vr.RecordStoreTime = v.GetInt64("RecordStoreTime")
	vr.RecordSendTime = v.GetInt64("RecordSendTime")

	vr.CitizenVehicleType = v.GetInt("CitizenVehicleType")
	vr.CitizenVehicleColor = v.GetInt("CitizenVehicleColor")
	vr.CitizenVehicleModel = string(v.GetStringBytes("CitizenVehicleModel"))
	vr.CitizenVehicleDistance = v.GetInt("CitizenVehicleDistance")
	vr.CitizenVehicleDegree = v.GetInt("CitizenVehicleDegree")
	vr.CitizenPlateNumber = string(v.GetStringBytes("CitizenPlateNumber"))
	vr.CitizenVehiclePlateNumberType = v.GetInt("CitizenVehiclePlateNumberType")
	vr.CitizenVehiclePlateNumberColor = v.GetInt("CitizenVehiclePlateNumberColor")

	vr.OCRAccuracy = v.GetFloat64("OCRAccuracy")
	vr.IsCitizenVehicleDistorted = v.GetBool("IsCitizenVehicleDistorted")
	vr.IsCitizenVehiclePlateNumberVisible = v.GetBool("IsCitizenVehiclePlateNumberVisible")
	vr.Sent = v.GetBool("Sent")
	vr.Retries = v.GetInt("Retries")
	vr.BackoffTime = v.GetInt64("Retries")
	vr.CitizenParkType = v.GetInt("CitizenParkType")

	vr.RingID = v.GetInt64("RingId")
	vr.StreetID = v.GetInt64("StreetId")
	vr.SegmentID = v.GetInt64("SegmentId")
	vr.RoadCode = v.GetInt64("RoadCode")
	vr.IsJunction = v.GetBool("IsJunction")
	vr.ParkingLotID = v.GetInt64("ParkingLotId")

	vr.UserID = string(v.GetStringBytes("UserId"))
	vr.LPRVehicleID = string(v.GetStringBytes("LPRVehicleId"))
	vr.LPRSystemID = string(v.GetStringBytes("LPRSystemId"))
	vr.LPRSystemAppID = string(v.GetStringBytes("LPRSystemAppId"))
	vr.LPRSystemAppVersion = string(v.GetStringBytes("LPRSystemAppVersion"))

	vr.LPRVehicleGPSSpeed = v.GetFloat64("LPRVehicleGPSSpeed")
	vr.LPRVehicleIsGPSSignalValid = v.GetBool("LPRVehicleIsGPSSignalValid")
	vr.LPRVehicleGPSLatitude = v.GetFloat64("LPRVehicleGPSLatitude")
	vr.LPRVehicleGPSLongitude = v.GetFloat64("LPRVehicleGPSLongitude")
	vr.LPRVehicleGPSError = v.GetInt("LPRVehicleGPSError")

	vr.LPRVehicleRTKLatitude = v.GetFloat64("LPRVehicleRTKLatitude")
	vr.LPRVehicleRTKLongitude = v.GetFloat64("LPRVehicleRTKLongitude")
	vr.LPRVehicleRTKError = v.GetInt("LPRVehicleRTKError")

	if photosArray := v.GetArray("VehiclePhotos"); len(photosArray) > 0 {
		vr.VehiclePhotos = make([]*CitizenVehiclePhoto, len(photosArray))
		for i, photoJSON := range photosArray {
			var photo CitizenVehiclePhoto

			photo.RecordID = vr.RecordID
			photo.PhotoSequenceID = photoJSON.GetInt("PhotoSequenceId")
			photo.OCRAccuracy = photoJSON.GetFloat64("OCRAccuracy")
			photo.LPRVehicleCameraID = photoJSON.GetInt("LPRVehicleCameraId")
			photo.CitizenVehiclePhoto = string(photoJSON.GetStringBytes("CitizenVehiclePhoto"))
			photo.CitizenVehiclePhotoArea = string(photoJSON.GetStringBytes("CitizenVehiclePhotoArea"))
			photo.CitizenVehiclePlateCropPhoto = string(photoJSON.GetStringBytes("CitizenVehiclePlateCropPhoto"))
			photo.CitizenVehiclePhotoCaptureTime = photoJSON.GetInt64("CitizenVehiclePhotoCaptureTime")
			vr.VehiclePhotos[i] = &photo
		}
	}
	if len(vr.VehiclePhotos) < 1 {
		if photosArray := v.GetArray("CitizenVehiclePhotosArray"); len(photosArray) > 0 {
			vr.VehiclePhotos = make([]*CitizenVehiclePhoto, len(photosArray))
			for i, photoJSON := range photosArray {
				var photo CitizenVehiclePhoto

				photo.RecordID = vr.RecordID
				photo.PhotoSequenceID = photoJSON.GetInt("PhotoSequenceId")
				photo.OCRAccuracy = photoJSON.GetFloat64("OCRAccuracy")
				photo.LPRVehicleCameraID = photoJSON.GetInt("LPRVehicleCameraId")
				photo.CitizenVehiclePhoto = string(photoJSON.GetStringBytes("CitizenVehiclePhoto"))
				photo.CitizenVehiclePhotoArea = string(photoJSON.GetStringBytes("CitizenVehiclePhotoArea"))
				photo.CitizenVehiclePlateCropPhoto = string(photoJSON.GetStringBytes("CitizenVehiclePlateCropPhoto"))
				photo.CitizenVehiclePhotoCaptureTime = photoJSON.GetInt64("CitizenVehiclePhotoCaptureTime")
				vr.VehiclePhotos[i] = &photo
			}
		}
	}

	return nil
}
