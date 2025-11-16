package seeder

import (
	"context"
	"farin/domain/entity"
	"farin/domain/repository"
	"github.com/bxcodec/faker/v4"
	"math/rand"
	"time"
)

func SeedCalenders(
	contracts []*entity.Contract,
	cr *repository.CalenderRepository,
) ([]*entity.Calender, error) {
	weekdays := []entity.Weekday{
		entity.Saturday, entity.Sunday, entity.Monday,
		entity.Tuesday, entity.Wednesday, entity.Thursday, entity.Friday,
	}
	workShifts := []entity.WorkShift{
		entity.Morning, entity.Afternoon, entity.Both,
	}

	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)
	var calenders []*entity.Calender

	// Generate calender entries for each contract
	for _, contract := range contracts {
		// Determine the date range for the contract
		startDate := contract.StartDate
		endDate := contract.EndDate

		// Generate calender entries within the contract period
		currentDate := startDate
		for currentDate.Before(endDate) {
			calender := &entity.Calender{
				ContractID:     contract.ID,
				ShamsiDate:     convertToShamsiDate(currentDate),
				WorkDate:       currentDate.Unix(),
				Weekday:        weekdays[currentDate.Weekday()],
				Year:           currentDate.Year(),
				IsHoliday:      isHoliday(currentDate),
				WorkShift:      workShifts[r.Intn(len(workShifts))],
				Description:    faker.Sentence(),
				WorkShiftStart: getWorkShiftStart(currentDate),
				WorkShiftEnd:   getWorkShiftEnd(currentDate),
			}

			// Create the calender entry
			createdCalender, err := cr.Create(context.Background(), calender)
			if err != nil {
				return nil, err
			}
			calenders = append(calenders, createdCalender)

			// Move to next day
			currentDate = currentDate.AddDate(0, 0, 1)
		}
	}

	return calenders, nil
}

// convertToShamsiDate converts Gregorian date to Shamsi (Persian) date
// Note: This is a placeholder implementation. You'll need to implement
// actual Shamsi date conversion logic
func convertToShamsiDate(t time.Time) string {
	// Implement actual Shamsi date conversion
	return t.Format("2006-01-02")
}

// isHoliday determines if a given date is a holiday
// This is a simple implementation. You should replace with actual holiday logic
func isHoliday(t time.Time) bool {
	// Example: weekends are holidays
	return t.Weekday() == time.Saturday || t.Weekday() == time.Sunday
}

// getWorkShiftStart generates a work shift start time
func getWorkShiftStart(t time.Time) int64 {
	// Default morning shift start at 8:00 AM
	morning := time.Date(t.Year(), t.Month(), t.Day(), 8, 0, 0, 0, t.Location())
	return morning.Unix()
}

// getWorkShiftEnd generates a work shift end time
func getWorkShiftEnd(t time.Time) int64 {
	// Default morning shift end at 4:00 PM
	afternoon := time.Date(t.Year(), t.Month(), t.Day(), 16, 0, 0, 0, t.Location())
	return afternoon.Unix()
}
