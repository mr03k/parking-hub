package entity

import (
	"encoding/json"
	"github.com/valyala/fastjson"
)

// Event represents the data structure for citizen vehicle records.
type Event struct {
	ID                                 string  `json:"id" bson:"_id,omitempty"`
	CitizenVehicleColor                int     `bson:"citizenVehicleColor,omitempty" json:"citizenVehicleColor"`
	LPRVehicleGPSError                 float64 `bson:"lprVehicleGPSError,omitempty" json:"lprVehicleGPSError"`
	LPRVehicleCameraID                 int     `bson:"lprVehicleCameraId,omitempty" json:"lprVehicleCameraId"`
	RecordID                           string  `bson:"recordId,omitempty" json:"recordId"`
	RecordSendTime                     int64   `bson:"recordSendTime,omitempty" json:"recordSendTime"`
	IsCitizenVehicleDistorted          bool    `bson:"isCitizenVehicleDistorted,omitempty" json:"isCitizenVehicleDistorted"`
	CitizenVehiclePhoto                []byte  `bson:"citizenVehiclePhoto,omitempty" json:"citizenVehiclePhoto"`
	RecordStoreTime                    int64   `bson:"recordStoreTime,omitempty" json:"recordStoreTime"`
	CitizenVehicleModel                string  `bson:"citizenVehicleModel,omitempty" json:"citizenVehicleModel"`
	CitizenVehicleDistance             int     `bson:"citizenVehicleDistance,omitempty" json:"citizenVehicleDistance"`
	CitizenParkType                    int     `bson:"citizenParkType,omitempty" json:"citizenParkType"`
	LPRSystemID                        string  `bson:"lprSystemId,omitempty" json:"lprSystemId"`
	LPRVehicleIsGPSSignalValid         bool    `bson:"lprVehicleIsGPSSignalValid,omitempty" json:"lprVehicleIsGPSSignalValid"`
	CitizenVehicleType                 int     `bson:"citizenVehicleType,omitempty" json:"citizenVehicleType"`
	CitizenVehiclePlateNumberColor     int     `bson:"citizenVehiclePlateNumberColor,omitempty" json:"citizenVehiclePlateNumberColor"`
	IsCitizenVehiclePlateNumberVisible bool    `bson:"isCitizenVehiclePlateNumberVisible,omitempty" json:"isCitizenVehiclePlateNumberVisible"`
	UserID                             string  `bson:"userId,omitempty" json:"userId"`
	CitizenVehiclePhotoCaptureTime     int64   `bson:"citizenVehiclePhotoCaptureTime,omitempty" json:"citizenVehiclePhotoCaptureTime"`
	CitizenPlateNumber                 string  `bson:"citizenPlateNumber,omitempty" json:"citizenPlateNumber"`
	LPRSystemAppID                     string  `bson:"lprSystemAppId,omitempty" json:"lprSystemAppId"`
	LPRSystemAppVersion                string  `bson:"lprSystemAppVersion,omitempty" json:"lprSystemAppVersion"`
	LPRSChannelID                      string  `bson:"lprsChaninId,omitempty" json:"lprsChaninId"`
	LPRVehicleRTKError                 float64 `bson:"lprVehicleRTKError,omitempty" json:"lprVehicleRTKError"`
	CitizenVehiclePhotoArea            string  `bson:"citizenVehiclePhotoArea,omitempty" json:"citizenVehiclePhotoArea"`
	OCRAccuracy                        float64 `bson:"ocrAccuracy,omitempty" json:"ocrAccuracy"`
	LPRVehicleGPSSpeed                 float64 `bson:"lprVehicleGPSSpeed,omitempty" json:"lprVehicleGPSSpeed"`
	CitizenVehicleDegree               int     `bson:"citizenVehicleDegree,omitempty" json:"citizenVehicleDegree"`
	CitizenVehiclePlateNumberType      int     `bson:"citizenVehiclePlateNumberType,omitempty" json:"citizenVehiclePlateNumberType"`
	RingID                             string  `bson:"ringId,omitempty" json:"ringId"`
	StreetID                           string  `bson:"streetId,omitempty" json:"streetId"`
	LPRVehicleID                       string  `bson:"lprVehicleId,omitempty" json:"lprVehicleId"`
	//LPRVehicleGPSLocation              Geo     `bson:"lprVehicleGPSLocation,omitempty" json:"lprVehicleGPSLatitude"`
	//LPRVehicleRTKLocation              Geo     `bson:"lprVehicleRTKLocation,omitempty" json:"lprVehicleRTKLongitude"`
	Published string `bson:"published,omitempty" json:"published"`
}

// MarshalJSON marshals the Event struct to JSON.
func (cv *Event) MarshalJSON() ([]byte, error) {
	return json.Marshal(cv)
}

func (cv *Event) UnmarshalNonFileFieldsJSON(data []byte) error {
	var p fastjson.Parser
	v, err := p.ParseBytes(data)
	if err != nil {
		return err
	}

	cv.ID = string(v.GetStringBytes("id"))
	cv.CitizenVehicleColor = v.GetInt("citizenVehicleColor")
	cv.LPRVehicleGPSError = v.GetFloat64("lprVehicleGPSError")
	cv.LPRVehicleCameraID = v.GetInt("lprVehicleCameraId")
	cv.RecordID = string(v.GetStringBytes("recordId"))
	cv.RecordSendTime = v.GetInt64("recordSendTime")
	cv.IsCitizenVehicleDistorted = v.GetBool("isCitizenVehicleDistorted")
	cv.RecordStoreTime = v.GetInt64("recordStoreTime")
	cv.CitizenVehicleModel = string(v.GetStringBytes("citizenVehicleModel"))
	cv.CitizenVehicleDistance = v.GetInt("citizenVehicleDistance")
	cv.CitizenParkType = v.GetInt("citizenParkType")
	cv.LPRSystemID = string(v.GetStringBytes("lprSystemId"))
	cv.LPRVehicleIsGPSSignalValid = v.GetBool("lprVehicleIsGPSSignalValid")
	cv.CitizenVehicleType = v.GetInt("citizenVehicleType")
	cv.CitizenVehiclePlateNumberColor = v.GetInt("citizenVehiclePlateNumberColor")
	cv.IsCitizenVehiclePlateNumberVisible = v.GetBool("isCitizenVehiclePlateNumberVisible")
	cv.UserID = string(v.GetStringBytes("userId"))
	cv.CitizenVehiclePhotoCaptureTime = v.GetInt64("citizenVehiclePhotoCaptureTime")
	cv.CitizenPlateNumber = string(v.GetStringBytes("citizenPlateNumber"))
	cv.LPRSystemAppID = string(v.GetStringBytes("lprSystemAppId"))
	cv.LPRSystemAppVersion = string(v.GetStringBytes("lprSystemAppVersion"))
	cv.LPRSChannelID = string(v.GetStringBytes("lprsChaninId"))
	cv.LPRVehicleRTKError = v.GetFloat64("lprVehicleRTKError")
	cv.OCRAccuracy = v.GetFloat64("ocrAccuracy")
	cv.LPRVehicleGPSSpeed = v.GetFloat64("lprVehicleGPSSpeed")
	cv.CitizenVehicleDegree = v.GetInt("citizenVehicleDegree")
	cv.CitizenVehiclePlateNumberType = v.GetInt("citizenVehiclePlateNumberType")
	cv.RingID = string(v.GetStringBytes("ringId"))
	cv.StreetID = string(v.GetStringBytes("streetId"))
	cv.LPRVehicleID = string(v.GetStringBytes("lprVehicleId"))
	cv.Published = string(v.GetStringBytes("published"))
	//cv.LPRVehicleGPSLocation = Geo{
	//	Coordinates: []float64{v.GetFloat64("lprVehicleGPSLocation", "latitude"),
	//		v.GetFloat64("lprVehicleGPSLocation", "longitude")},
	//	Type: "Point",
	//}
	//cv.LPRVehicleRTKLocation = Geo{
	//	Coordinates: []float64{v.GetFloat64("lprVehicleRTKLocation", "latitude"),
	//		v.GetFloat64("lprVehicleRTKLocation", "longitude")},
	//	Type: "Point",
	//}
	cv.CitizenVehiclePhoto = v.GetStringBytes("citizenVehiclePhoto")

	return nil
}
