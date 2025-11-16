package repo

import (
	"application/internal/entity"
	"errors"
	"github.com/google/uuid"
	"sync"
	"time"
)

// DistrictRepository manages in-memory storage for districts.
type DistrictRepository struct {
	data  map[string]entity.District // In-memory storage for districts
	mutex sync.RWMutex               // Mutex for thread-safe operations
}

// NewDistrictRepository creates a new instance of DistrictRepository with seeded data.
func NewDistrictRepository() *DistrictRepository {
	repo := &DistrictRepository{
		data: make(map[string]entity.District),
	}

	// Seed the repository with initial data
	repo.seedData()

	return repo
}

// seedData initializes the repository with mock districts.
func (r *DistrictRepository) seedData() {
	now := time.Now()

	districts := []entity.District{
		{
			ID:           uuid.New(),
			DistrictName: "District 5",
			DistrictCode: "D5",
			CityID:       uuid.New(),
			GeoBoundary:  "POLYGON((...))",
			Population:   100000,
			Area:         150.5,
			CreatedAt:    now,
		},
		{
			ID:           uuid.New(),
			DistrictName: "District 10",
			DistrictCode: "D10",
			CityID:       uuid.New(),
			GeoBoundary:  "POLYGON((...))",
			Population:   200000,
			Area:         250.3,
			CreatedAt:    now,
		},
	}

	// Add seeded districts to the data map
	for _, district := range districts {
		r.data[district.ID.String()] = district
	}
}

// GetDistricts retrieves all districts.
func (r *DistrictRepository) GetDistricts() []entity.District {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	districts := make([]entity.District, 0, len(r.data))
	for _, district := range r.data {
		districts = append(districts, district)
	}

	return districts
}

// GetDistrictByID retrieves a district by its ID.
func (r *DistrictRepository) GetDistrictByID(id string) (*entity.District, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	district, exists := r.data[id]
	if !exists {
		return nil, errors.New("district not found")
	}

	return &district, nil
}

// AddDistrict adds a new district to the repository.
func (r *DistrictRepository) AddDistrict(district entity.District) string {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	district.ID = uuid.New() // Generate a new UserID for the district
	district.CreatedAt = time.Now()

	r.data[district.ID.String()] = district
	return district.ID.String()
}

// UpdateDistrict updates an existing district's details.
func (r *DistrictRepository) UpdateDistrict(id string, updatedDistrict entity.District) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	district, exists := r.data[id]
	if !exists {
		return errors.New("district not found")
	}

	updatedDistrict.ID = district.ID
	updatedDistrict.CreatedAt = district.CreatedAt // Preserve original creation timestamp

	r.data[id] = updatedDistrict
	return nil
}

// DeleteDistrict removes a district from the repository by its ID.
func (r *DistrictRepository) DeleteDistrict(id string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	_, exists := r.data[id]
	if !exists {
		return errors.New("district not found")
	}

	delete(r.data, id)
	return nil
}
