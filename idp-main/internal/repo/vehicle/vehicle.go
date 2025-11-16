package vehicle

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"application/internal/biz/device"
	"application/internal/datasource"
	device2 "application/internal/entity/device"
)

// vehicleRepo
type vehicleRepo struct {
	logger *slog.Logger
	db     *sql.DB
}

// NewVehicleRepo
func NewVehicleRepo(logger *slog.Logger, ds *datasource.Datasource) device.VehicleRepositoryInterface {
	return &vehicleRepo{
		logger: logger.With("layer", "VehicleRepo"),
		db:     ds.DBpsql,
	}
}

// CreateVehicle
func (r *vehicleRepo) CreateVehicle(ctx context.Context, vehicle *device2.Vehicle) (*device2.Vehicle, error) {
	logger := r.logger.With("method", "CreateVehicle")
	logger.Debug("repository CreateVehicle")

	contractorID := sql.NullInt64{
		Int64: vehicle.ContractorID,
	}
	// Insert vehicle and retrieve ID
	err := r.db.QueryRowContext(
		ctx,
		`INSERT INTO vehicles (code_vehicle, vin, plate_license, type_vehicle, brand, model, color, manufacture_of_year, kilometers_initial, expiry_insurance_party_third, expiry_insurance_body, image_document_vehicle, image_card_vehicle, third_party_insurance_image, body_insurance_image, id_contractor, status, description)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18)
		 RETURNING id`,
		vehicle.CodeVehicle, vehicle.VIN, vehicle.PlateLicense, vehicle.TypeVehicle, vehicle.Brand, vehicle.Model, vehicle.Color, vehicle.ManufactureOfYear, vehicle.KilometersInitial, vehicle.ExpiryInsurancePartyThird, vehicle.ExpiryInsuranceBody, vehicle.ImageDocumentVehicle, vehicle.ImageCardVehicle,
		vehicle.ThirdPartyInsuranceImage, vehicle.BodyInsuranceImage, contractorID, vehicle.Status, vehicle.Description,
	).Scan(&vehicle.ID)
	if err != nil {
		logger.Error("error CreateVehicle", "error", err)
		return nil, err
	}

	return vehicle, nil
}

// GetVehicle
func (r *vehicleRepo) GetVehicle(ctx context.Context, vehicleID string) (*device2.Vehicle, error) {
	logger := r.logger.With("method", "GetVehicle")
	logger.Debug("repository GetVehicle")

	vehicle := device2.Vehicle{}
	contractorID := sql.NullInt64{}

	err := r.db.QueryRowContext(
		ctx,
		`SELECT id, code_vehicle, vin, plate_license, type_vehicle, brand, model, color, manufacture_of_year, kilometers_initial, expiry_insurance_party_third, expiry_insurance_body, image_document_vehicle, image_card_vehicle, third_party_insurance_image, body_insurance_image, id_contractor, status, description
		 FROM vehicles WHERE id = $1`,
		vehicleID,
	).Scan(&vehicle.ID, &vehicle.CodeVehicle, &vehicle.VIN, &vehicle.PlateLicense, &vehicle.TypeVehicle,
		&vehicle.Brand, &vehicle.Model, &vehicle.Color, &vehicle.ManufactureOfYear, &vehicle.KilometersInitial,
		&vehicle.ExpiryInsurancePartyThird, &vehicle.ExpiryInsuranceBody, &vehicle.ImageDocumentVehicle,
		&vehicle.ImageCardVehicle, &vehicle.ThirdPartyInsuranceImage, &vehicle.BodyInsuranceImage,
		&contractorID, &vehicle.Status, &vehicle.Description)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, device.ErrVehicleNotFound
	}
	if err != nil {
		logger.Error("error GetVehicle", "error", err)
		return nil, err
	}

	vehicle.ContractorID = contractorID.Int64

	return &vehicle, nil
}

// ListVehicles
func (r *vehicleRepo) ListVehicles(ctx context.Context) ([]device2.Vehicle, error) {
	logger := r.logger.With("method", "ListVehicles")
	logger.Debug("repository ListVehicles")

	vehicles := []device2.Vehicle{}
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT id, code_vehicle, vin, plate_license, type_vehicle, brand, model, color, manufacture_of_year,
       kilometers_initial, expiry_insurance_party_third, expiry_insurance_body, image_document_vehicle,
       image_card_vehicle, third_party_insurance_image, body_insurance_image, id_contractor, status, description
		 FROM vehicles`,
	)
	if err != nil {
		logger.Error("error ListVehicles", "error", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		vehicle := device2.Vehicle{}
		contractorID := sql.NullInt64{}
		err := rows.Scan(&vehicle.ID, &vehicle.CodeVehicle, &vehicle.VIN, &vehicle.PlateLicense, &vehicle.TypeVehicle,
			&vehicle.Brand, &vehicle.Model, &vehicle.Color, &vehicle.ManufactureOfYear, &vehicle.KilometersInitial,
			&vehicle.ExpiryInsurancePartyThird, &vehicle.ExpiryInsuranceBody, &vehicle.ImageDocumentVehicle,
			&vehicle.ImageCardVehicle, &vehicle.ThirdPartyInsuranceImage, &vehicle.BodyInsuranceImage,
			&contractorID, &vehicle.Status, &vehicle.Description)
		if err != nil {
			logger.Error("error scanning row in ListVehicles", "error", err)
			return nil, err
		}
		vehicle.ContractorID = contractorID.Int64
		vehicles = append(vehicles, vehicle)
	}

	return vehicles, nil
}

// DeleteVehicle
func (r *vehicleRepo) DeleteVehicle(ctx context.Context, vehicleID string) error {
	logger := r.logger.With("method", "DeleteVehicle")
	logger.Debug("repository DeleteVehicle")

	res, err := r.db.ExecContext(
		ctx,
		"DELETE FROM vehicles WHERE id = $1",
		vehicleID,
	)
	if err != nil {
		logger.Error("error DeleteVehicle", "error", err)
		return err
	}
	n, err := res.RowsAffected()
	if err == nil && n < 1 {
		return device.ErrVehicleNotFound
	}

	return nil
}
