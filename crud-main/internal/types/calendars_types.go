package types

import (
	"time"

	"github.com/zhufuyi/sponge/pkg/ggorm/query"
)

var _ time.Time

// Tip: suggested filling in the binding rules https://github.com/go-playground/validator in request struct fields tag.

// CreateCalendarsRequest request params
type CreateCalendarsRequest struct {
	IDContract  string    `json:"iDContract" binding:""`
	ShamsiDate  string    `json:"shamsiDate" binding:""`
	WorkDate    time.Time `json:"workDate" binding:""`
	Weekday     string    `json:"weekday" binding:""`
	Year        int       `json:"year" binding:""`
	HolidayIs   bool      `json:"holidayIs" binding:""`
	ShiftWork   string    `json:"shiftWork" binding:""`
	Description string    `json:"description" binding:""`
}

// UpdateCalendarsByIDRequest request params
type UpdateCalendarsByIDRequest struct {
	ID uint64 `json:"id" binding:""` // uint64 id

	IDContract  string    `json:"iDContract" binding:""`
	ShamsiDate  string    `json:"shamsiDate" binding:""`
	WorkDate    time.Time `json:"workDate" binding:""`
	Weekday     string    `json:"weekday" binding:""`
	Year        int       `json:"year" binding:""`
	HolidayIs   bool      `json:"holidayIs" binding:""`
	ShiftWork   string    `json:"shiftWork" binding:""`
	Description string    `json:"description" binding:""`
}

// CalendarsObjDetail detail
type CalendarsObjDetail struct {
	ID uint64 `json:"id"` // convert to uint64 id

	IDContract  string    `json:"iDContract"`
	ShamsiDate  string    `json:"shamsiDate"`
	WorkDate    time.Time `json:"workDate"`
	Weekday     string    `json:"weekday"`
	Year        int       `json:"year"`
	HolidayIs   bool      `json:"holidayIs"`
	ShiftWork   string    `json:"shiftWork"`
	Description string    `json:"description"`
	CreatedAt   int       `json:"createdAt"`
}

// CreateCalendarsReply only for api docs
type CreateCalendarsReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		ID uint64 `json:"id"` // id
	} `json:"data"` // return data
}

// DeleteCalendarsByIDReply only for api docs
type DeleteCalendarsByIDReply struct {
	Result
}

// UpdateCalendarsByIDReply only for api docs
type UpdateCalendarsByIDReply struct {
	Result
}

// GetCalendarsByIDReply only for api docs
type GetCalendarsByIDReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		Calendars CalendarsObjDetail `json:"calendars"`
	} `json:"data"` // return data
}

// ListCalendarssRequest request params
type ListCalendarssRequest struct {
	query.Params
}

// ListCalendarssReply only for api docs
type ListCalendarssReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		Calendarss []CalendarsObjDetail `json:"calendarss"`
	} `json:"data"` // return data
}
