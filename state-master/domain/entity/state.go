package entity

import "github.com/valyala/fastjson"

type State struct {
	Base
	UserID                 string  `gorm:"type:uuid;null"`
	RecordID               string  `gorm:"type:uuid;not null;index"`
	LPRSystemID            string  `gorm:"type:varchar(255);not null"`
	LPRVehicleID           string  `gorm:"type:uuid;not null"`
	LPRSystemAppID         string  `gorm:"type:uuid;not null"`
	LPRSystemAppVersion    string  `gorm:"type:varchar(20);not null"`
	LPRVehicleGPSLatitude  float64 `gorm:"type:double precision;not null"`
	LPRVehicleGPSLongitude float64 `gorm:"type:double precision;not null"`
	LPRVehicleGPSSpeed     float32 `gorm:"type:float;not null"`
	LPRVehicleGPSError     int     `gorm:"type:integer;not null"`
	RecordStoreTime        int64   `gorm:"not null;index"`
	RecordSendTime         int64   `gorm:"not null;index"`
	ServerAvailability     bool    `gorm:"not null"`
	ServerPingTime         int     `gorm:"type:integer;not null"`
}

func (s *State) UnmarshalJSON(data []byte) error {
	var p fastjson.Parser
	v, err := p.ParseBytes(data)
	if err != nil {
		return err
	}

	// User and Record Identification
	if userIDBytes := v.GetStringBytes("UserId"); len(userIDBytes) > 0 {
		s.UserID = string(userIDBytes)
	}

	if recordIDBytes := v.GetStringBytes("RecordId"); len(recordIDBytes) > 0 {
		s.RecordID = string(recordIDBytes)
	}

	// LPR System and Vehicle data
	if lprSystemIDBytes := v.GetStringBytes("LPRSystemId"); len(lprSystemIDBytes) > 0 {
		s.LPRSystemID = string(lprSystemIDBytes)
	}

	if lprVehicleIDBytes := v.GetStringBytes("LPRVehicleId"); len(lprVehicleIDBytes) > 0 {
		s.LPRVehicleID = string(lprVehicleIDBytes)
	}

	if lprSystemAppIDBytes := v.GetStringBytes("LPRSystemAppId"); len(lprSystemAppIDBytes) > 0 {
		s.LPRSystemAppID = string(lprSystemAppIDBytes)
	}

	if lprSystemAppVersionBytes := v.GetStringBytes("LPRSystemAppVersion"); len(lprSystemAppVersionBytes) > 0 {
		s.LPRSystemAppVersion = string(lprSystemAppVersionBytes)
	}

	// GPS data
	s.LPRVehicleGPSLatitude = v.GetFloat64("LPRVehicleGPSLatitude")
	s.LPRVehicleGPSLongitude = v.GetFloat64("LPRVehicleGPSLongitude")
	s.LPRVehicleGPSSpeed = float32(v.GetFloat64("LPRVehicleGPSSpeed"))
	s.LPRVehicleGPSError = v.GetInt("LPRVehicleGPSError")

	// Timing data
	s.RecordStoreTime = v.GetInt64("RecordStoreTime")
	s.RecordSendTime = v.GetInt64("RecordSendTime")

	// Server data
	s.ServerAvailability = v.GetBool("ServerAvailability")
	s.ServerPingTime = v.GetInt("ServerPingTime")

	return nil
}
