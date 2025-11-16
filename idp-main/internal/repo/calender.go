package repo

import (
	"application/internal/entity"
	"errors"
	"github.com/google/uuid"
	"sync"
	"time"
)

// WorkCalendarRepository manages in-memory storage for work calendars.
type WorkCalendarRepository struct {
	data  map[string]entity.WorkCalendar // In-memory storage for work calendar records
	mutex sync.RWMutex                   // Mutex for thread-safe operations
}

// NewWorkCalendarRepository creates a new instance of WorkCalendarRepository with seeded data.
func NewWorkCalendarRepository() *WorkCalendarRepository {
	repo := &WorkCalendarRepository{
		data: make(map[string]entity.WorkCalendar),
	}

	// Seed the repository with initial data
	repo.seedData()

	return repo
}

// seedData initializes the repository with mock work calendar records.
func (r *WorkCalendarRepository) seedData() {
	now := time.Now()

	workCalendars := []entity.WorkCalendar{
		{
			CalendarID:     uuid.New(),
			ContractID:     uuid.New(),
			ShamsiDate:     "01/01/1402",
			WorkDate:       now.AddDate(-1, 0, 0), // A day one year ago
			Weekday:        "Saturday",
			Year:           1402,
			IsHoliday:      false,
			WorkShift:      "Morning",
			Description:    "First working day of the Shamsi year 1402",
			CreatedAt:      now.AddDate(-1, 0, 0),
			UpdatedAt:      now.AddDate(-1, 0, 0),
			WorkShiftStart: time.Date(0, 1, 1, 8, 0, 0, 0, time.UTC),
			WorkShiftEnd:   time.Date(0, 1, 1, 16, 0, 0, 0, time.UTC),
		},
		{
			CalendarID:  uuid.New(),
			ContractID:  uuid.New(),
			ShamsiDate:  "02/01/1402",
			WorkDate:    now.AddDate(-1, 0, 1), // A day after the first record
			Weekday:     "Sunday",
			Year:        1402,
			IsHoliday:   true,
			WorkShift:   "Both",
			Description: "Holiday for celebrations",
			CreatedAt:   now.AddDate(-1, 0, 1),
			UpdatedAt:   now.AddDate(-1, 0, 1),
		},
	}

	// Add seeded work calendars to the data map
	for _, calendar := range workCalendars {
		r.data[calendar.CalendarID.String()] = calendar
	}
}

// GetWorkCalendars retrieves all work calendar records.
func (r *WorkCalendarRepository) GetWorkCalendars() []entity.WorkCalendar {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	workCalendars := make([]entity.WorkCalendar, 0, len(r.data))
	for _, calendar := range r.data {
		workCalendars = append(workCalendars, calendar)
	}

	return workCalendars
}

// GetWorkCalendarByID retrieves a work calendar record by its ID.
func (r *WorkCalendarRepository) GetWorkCalendarByID(id string) (*entity.WorkCalendar, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	calendar, exists := r.data[id]
	if !exists {
		return nil, errors.New("work calendar not found")
	}

	return &calendar, nil
}

// AddWorkCalendar adds a new work calendar record to the repository.
func (r *WorkCalendarRepository) AddWorkCalendar(calendar entity.WorkCalendar) string {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	calendar.CalendarID = uuid.New() // Generate a new UserID for the work calendar
	calendar.CreatedAt = time.Now()
	calendar.UpdatedAt = time.Now()

	r.data[calendar.CalendarID.String()] = calendar
	return calendar.CalendarID.String()
}

// UpdateWorkCalendar updates an existing work calendar record's details.
func (r *WorkCalendarRepository) UpdateWorkCalendar(id string, updatedCalendar entity.WorkCalendar) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	calendar, exists := r.data[id]
	if !exists {
		return errors.New("work calendar not found")
	}

	updatedCalendar.CalendarID = calendar.CalendarID
	updatedCalendar.CreatedAt = calendar.CreatedAt // Preserve original creation timestamp
	updatedCalendar.UpdatedAt = time.Now()

	r.data[id] = updatedCalendar
	return nil
}

// DeleteWorkCalendar removes a work calendar record from the repository by its ID.
func (r *WorkCalendarRepository) DeleteWorkCalendar(id string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	_, exists := r.data[id]
	if !exists {
		return errors.New("work calendar not found")
	}

	delete(r.data, id)
	return nil
}
