package model

import (
	"time"
)

type Rates struct {
	ID                     uint64     `gorm:"column:id;type:uuid;primary_key" json:"id"`
	Code                   string     `gorm:"column:code;type:varchar(20);NOT NULL" json:"code"`
	RoadCategoryID         string     `gorm:"column:road_category_id;type:uuid;NOT NULL" json:"roadCategoryID"`
	TimeCycleMinutes       int        `gorm:"column:time_cycle_minutes;type:int4;NOT NULL" json:"timeCycleMinutes"`
	RateMultiplier         string     `gorm:"column:rate_multiplier;type:numeric;NOT NULL" json:"rateMultiplier"`
	PeakHourMultiplier     string     `gorm:"column:peak_hour_multiplier;type:numeric" json:"peakHourMultiplier"`
	GoodPercentage         int        `gorm:"column:good_percentage;type:int4" json:"goodPercentage"`
	NormalSettlementPeriod int        `gorm:"column:normal_settlement_period;type:int4" json:"normalSettlementPeriod"`
	LatePenalty            string     `gorm:"column:late_penalty;type:numeric" json:"latePenalty"`
	LatePenaltyMax         string     `gorm:"column:late_penalty_max;type:numeric" json:"latePenaltyMax"`
	ValidFrom              *time.Time `gorm:"column:valid_from;type:date;NOT NULL" json:"validFrom"`
	ValidTo                *time.Time `gorm:"column:valid_to;type:date" json:"validTo"`
	Description            string     `gorm:"column:description;type:text" json:"description"`
	StartTime              string     `gorm:"column:start_time;type:varchar(20)" json:"startTime"`
	EndTime                string     `gorm:"column:end_time;type:varchar(20)" json:"endTime"`
	CityID                 string     `gorm:"column:city_id;type:uuid;NOT NULL" json:"cityID"`
	ApprovalNumber         string     `gorm:"column:approval_number;type:varchar(20)" json:"approvalNumber"`
	ApprovalDate           *time.Time `gorm:"column:approval_date;type:date" json:"approvalDate"`
	Year                   int        `gorm:"column:year;type:int4;NOT NULL" json:"year"`
	BaseRateID             string     `gorm:"column:base_rate_id;type:uuid" json:"baseRateID"`
	ExceptionsID           string     `gorm:"column:exceptions_id;type:uuid" json:"exceptionsID"`
	CreatedAt              int        `gorm:"column:created_at;type:int4;NOT NULL" json:"createdAt"`
}
