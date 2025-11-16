package dto

import "time"

// WorkCalendarResponse represents a summary of work calendar information.
type WorkCalendarResponse struct {
	ID          string    `json:"id"`
	ShamsiDate  string    `json:"shamsiDate"`
	WorkDate    time.Time `json:"workDate"`
	Weekday     string    `json:"weekday"`
	Year        int       `json:"year"`
	IsHoliday   bool      `json:"isHoliday"`
	WorkShift   string    `json:"workShift"`
	Description string    `json:"description,omitempty"`
}

// WorkCalendarDetailResponse represents detailed work calendar information.
type WorkCalendarDetailResponse struct {
	ID             string    `json:"id"`
	ContractID     string    `json:"contractId"`
	ShamsiDate     string    `json:"shamsiDate"`
	WorkDate       time.Time `json:"workDate"`
	Weekday        string    `json:"weekday"`
	Year           int       `json:"year"`
	IsHoliday      bool      `json:"isHoliday"`
	WorkShift      string    `json:"workShift"`
	Description    string    `json:"description,omitempty"`
	WorkShiftStart time.Time `json:"workShiftStart,omitempty"`
	WorkShiftEnd   time.Time `json:"workShiftEnd,omitempty"`
}

// WorkCalendarListResponse represents the response for a list of work calendar records.
type WorkCalendarListResponse struct {
	WorkCalendars []WorkCalendarResponse `json:"workCalendars"`
}
