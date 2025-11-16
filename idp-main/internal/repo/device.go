package repo

import (
	"application/internal/entity/device"
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
)

// DeviceRepository manages in-memory storage for license plate reader devices.
type DeviceRepository struct {
	data  map[string]device.LicensePlateReaderDevice // In-memory storage for devices
	mutex sync.RWMutex                               // Mutex for thread-safe operations
}

// NewDeviceRepository creates a new instance of DeviceRepository with seeded data.
func NewDeviceRepository() *DeviceRepository {
	repo := &DeviceRepository{
		data: make(map[string]device.LicensePlateReaderDevice),
	}

	// Seed the repository with initial data
	repo.seedData()

	return repo
}

// seedData initializes the repository with mock devices.
func (r *DeviceRepository) seedData() {
	now := time.Now()

	devices := []device.LicensePlateReaderDevice{
		{
			ID:                  uuid.New(),
			CodeDevice:          "D12345",
			NumberSerial:        "SN987654321",
			Model:               "ModelX",
			DateInstallation:    now.AddDate(-1, 0, 0), // Installed 1 year ago
			DateExpiryWarranty:  now.AddDate(1, 0, 0),  // Warranty expires in 1 year
			DateExpiryInsurance: now.AddDate(0, 6, 0),  // Insurance expires in 6 months
			ClassDevice:         "Premium",
			ImageContractURL:    "https://example.com/contract/device1.jpg",
			ImageInsuranceURL:   "https://example.com/insurance/device1.jpg",
			ContractorID:        uuid.New(),
			Description:         "First seeded license plate reader device",
			CreatedAt:           now,
		},
		{
			ID:                  uuid.New(),
			CodeDevice:          "D67890",
			NumberSerial:        "SN123456789",
			Model:               "ModelY",
			DateInstallation:    now.AddDate(-2, 0, 0), // Installed 2 years ago
			DateExpiryWarranty:  now,                   // Warranty expired today
			DateExpiryInsurance: now.AddDate(1, 0, 0),  // Insurance expires in 1 year
			ClassDevice:         "Standard",
			ImageContractURL:    "https://example.com/contract/device2.jpg",
			ImageInsuranceURL:   "https://example.com/insurance/device2.jpg",
			ContractorID:        uuid.New(),
			Description:         "Second seeded license plate reader device",
			CreatedAt:           now,
		},
	}

	// Add seeded devices to the data map
	for _, device := range devices {
		r.data[device.ID.String()] = device
	}
}

// GetDevices retrieves all devices.
func (r *DeviceRepository) GetDevices() []device.LicensePlateReaderDevice {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	devices := make([]device.LicensePlateReaderDevice, 0, len(r.data))
	for _, device := range r.data {
		devices = append(devices, device)
	}

	return devices
}

// GetDeviceByID retrieves a device by its ID.
func (r *DeviceRepository) GetDeviceByID(id string) (*device.LicensePlateReaderDevice, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	device, exists := r.data[id]
	if !exists {
		return nil, errors.New("device not found")
	}

	return &device, nil
}

// AddDevice adds a new device to the repository.
func (r *DeviceRepository) AddDevice(device device.LicensePlateReaderDevice) string {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	device.ID = uuid.New() // Generate a new UserID for the device
	device.CreatedAt = time.Now()

	r.data[device.ID.String()] = device
	return device.ID.String()
}

// UpdateDevice updates an existing device's details.
func (r *DeviceRepository) UpdateDevice(id string, updatedDevice device.LicensePlateReaderDevice) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	device, exists := r.data[id]
	if !exists {
		return errors.New("device not found")
	}

	updatedDevice.ID = device.ID
	updatedDevice.CreatedAt = device.CreatedAt       // Preserve original creation timestamp
	updatedDevice.ContractorID = device.ContractorID // Preserve ContractorID
	r.data[id] = updatedDevice
	return nil
}

// DeleteDevice removes a device from the repository by its ID.
func (r *DeviceRepository) DeleteDevice(id string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	_, exists := r.data[id]
	if !exists {
		return errors.New("device not found")
	}

	delete(r.data, id)
	return nil
}
