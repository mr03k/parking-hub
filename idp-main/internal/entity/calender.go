package entity

import (
	"time"

	"github.com/google/uuid"
)

// WorkCalendar represents the entity for a work calendar record.
type WorkCalendar struct {
	CalendarID     uuid.UUID `json:"calendarId"`               // Unique identifier for the calendar
	ContractID     uuid.UUID `json:"contractId"`               // Foreign key to the contracts table
	ShamsiDate     string    `json:"shamsiDate"`               // Shamsi date in string format (e.g., "02/05/1402")
	WorkDate       time.Time `json:"workDate"`                 // Gregorian date for the work day
	Weekday        string    `json:"weekday"`                  // Weekday name (e.g., "Saturday", "Sunday")
	Year           int       `json:"year"`                     // Year of the Shamsi date
	IsHoliday      bool      `json:"isHoliday"`                // Whether the day is a holiday
	WorkShift      string    `json:"workShift"`                // Work shift type (e.g., "Morning", "Afternoon", "Both")
	Description    string    `json:"description,omitempty"`    // Additional description for the day
	CreatedAt      time.Time `json:"createdAt"`                // Timestamp when the calendar record was created
	UpdatedAt      time.Time `json:"updatedAt"`                // Timestamp when the calendar record was last updated
	WorkShiftStart time.Time `json:"workShiftStart,omitempty"` // Start time of the work shift
	WorkShiftEnd   time.Time `json:"workShiftEnd,omitempty"`   // End time of the work shift
}
