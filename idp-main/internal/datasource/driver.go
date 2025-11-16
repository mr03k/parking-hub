package datasource

import (
	"application/internal/entity"
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
)

// DriverRepository manages in-memory storage for drivers.
type DriverRepository struct {
	data  map[string]entity.Driver // In-memory storage for drivers
	mutex sync.RWMutex             // Mutex for thread-safe operations
}

// NewDriverRepository creates a new instance of DriverRepository with seeded data.
func NewDriverRepository() *DriverRepository {
	repo := &DriverRepository{
		data: make(map[string]entity.Driver),
	}

	// Seed the repository with initial data
	repo.seedData()

	return repo
}

// seedData initializes the repository with mock drivers.
func (r *DriverRepository) seedData() {
	now := time.Now()

	drivers := []entity.Driver{
		{
			ID:                  uuid.NewString(),
			Address:             "Tehran, Iran Khodro, Karami St, Tolai Alley",
			DriverType:          "Primary",
			ShiftType:           "Morning",
			EmploymentStatus:    "Active",
			EmploymentStartDate: &now,
			Description:         "Seeded driver 1",
			DriverPhotoURL:      "https://example.com/photos/driver1.jpg",
			IDCardImageURL:      "https://example.com/images/idcard1.jpg",
			CreatedAt:           now,
			UpdatedAt:           now,
		},
		{
			ID:                  uuid.NewString(),
			Address:             "Mashhad, Vakilabad Blvd",
			DriverType:          "Reserve",
			ShiftType:           "Evening",
			EmploymentStatus:    "Inactive",
			EmploymentStartDate: &now,
			Description:         "Seeded driver 2",
			DriverPhotoURL:      "https://example.com/photos/driver2.jpg",
			IDCardImageURL:      "https://example.com/images/idcard2.jpg",
			CreatedAt:           now,
			UpdatedAt:           now,
		},
	}

	// Add seeded drivers to the data map
	for _, driver := range drivers {
		r.data[driver.ID] = driver
	}
}

// GetDrivers retrieves all drivers.
func (r *DriverRepository) GetDrivers() []entity.Driver {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	drivers := make([]entity.Driver, 0, len(r.data))
	for _, driver := range r.data {
		drivers = append(drivers, driver)
	}

	return drivers
}

// GetDriverByID retrieves a driver by their ID.
func (r *DriverRepository) GetDriverByID(id string) (*entity.Driver, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	driver, exists := r.data[id]
	if !exists {
		return nil, errors.New("driver not found")
	}

	return &driver, nil
}

// AddDriver adds a new driver to the repository.
func (r *DriverRepository) AddDriver(driver entity.Driver) string {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	driver.ID = uuid.NewString() // Generate a new UserID for the driver
	driver.CreatedAt = time.Now()
	driver.UpdatedAt = time.Now()

	r.data[driver.ID] = driver
	return driver.ID
}

// UpdateDriver updates an existing driver's details.
func (r *DriverRepository) UpdateDriver(id string, updatedDriver entity.Driver) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	driver, exists := r.data[id]
	if !exists {
		return errors.New("driver not found")
	}

	updatedDriver.ID = id
	updatedDriver.CreatedAt = driver.CreatedAt // Preserve original creation timestamp
	updatedDriver.UpdatedAt = time.Now()

	r.data[id] = updatedDriver
	return nil
}

// DeleteDriver removes a driver from the repository by their ID.
func (r *DriverRepository) DeleteDriver(id string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	_, exists := r.data[id]
	if !exists {
		return errors.New("driver not found")
	}

	delete(r.data, id)
	return nil
}
