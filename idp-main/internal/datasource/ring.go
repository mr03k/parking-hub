package datasource

import (
	"application/internal/entity"
	"errors"
	"github.com/google/uuid"
	"sync"
	"time"
)

// RingRepository manages in-memory storage for rings.
type RingRepository struct {
	data  map[string]entity.Ring // In-memory store
	mutex sync.RWMutex           // Mutex for thread-safe access
}

// NewRingRepository creates a new instance of RingRepository.
func NewRingRepository() *RingRepository {
	repo := &RingRepository{
		data: make(map[string]entity.Ring),
	}

	// Seed the repository with initial data
	repo.seedData()

	return repo
}

// seedData initializes the repository with mock data.
func (r *RingRepository) seedData() {
	id1 := uuid.NewString()
	r.data[id1] = entity.Ring{
		ID:                   id1,
		RingName:             "Ring 1",
		RingCode:             "R001",
		RingLength:           5.2,
		RingBoundary:         "POLYGON((...))",
		ParkingSpots:         507,
		DisabledParkingSpots: 7,
		TrafficSigns:         79,
		DisabledTrafficSigns: 7,
		StartPoint:           "POINT(34.0522 -118.2437)",
		BufferDistance:       100.0,
		Description:          "This is a test record for Ring 1",
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
	}

	id2 := uuid.NewString()
	r.data[id2] = entity.Ring{
		ID:                   id2,
		RingName:             "Ring 2",
		RingCode:             "R002",
		RingLength:           8.5,
		RingBoundary:         "POLYGON((...))",
		ParkingSpots:         320,
		DisabledParkingSpots: 10,
		TrafficSigns:         55,
		DisabledTrafficSigns: 5,
		StartPoint:           "POINT(40.7128 -74.0060)",
		BufferDistance:       200.0,
		Description:          "This is a test record for Ring 2",
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
	}
}

// GetRings returns a list of all rings.
func (r *RingRepository) GetRings() []entity.Ring {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	rings := make([]entity.Ring, 0, len(r.data))
	for _, ring := range r.data {
		rings = append(rings, ring)
	}

	return rings
}

// GetRingDetail returns the details of a single ring by its ID.
func (r *RingRepository) GetRingDetail(ringID string) (*entity.Ring, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	ring, exists := r.data[ringID]
	if !exists {
		return nil, errors.New("ring not found")
	}

	return &ring, nil
}

// AddRing adds a new ring to the repository.
func (r *RingRepository) AddRing(ring entity.Ring) string {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	ring.ID = uuid.NewString() // Generate a new ID
	ring.CreatedAt = time.Now()
	ring.UpdatedAt = time.Now()

	r.data[ring.ID] = ring
	return ring.ID
}

// UpdateRing updates an existing ring in the repository.
func (r *RingRepository) UpdateRing(ringID string, updatedRing entity.Ring) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	_, exists := r.data[ringID]
	if !exists {
		return errors.New("ring not found")
	}

	updatedRing.ID = ringID
	updatedRing.UpdatedAt = time.Now()
	r.data[ringID] = updatedRing
	return nil
}

// DeleteRing removes a ring from the repository by ID.
func (r *RingRepository) DeleteRing(ringID string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	_, exists := r.data[ringID]
	if !exists {
		return errors.New("ring not found")
	}

	delete(r.data, ringID)
	return nil
}
