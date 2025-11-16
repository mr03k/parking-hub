package seeder

import (
	"context"
	"farin/domain/entity"
	"farin/domain/repository"
	"fmt"
	"github.com/bxcodec/faker/v4"
	"math/rand"
	"time"
)

func SeedContracts(
	contractors []*entity.Contractor,
	cr *repository.ContractRepository,
) ([]*entity.Contract, error) {
	// Define possible contract types
	contractTypes := []string{
		"Service",
		"Maintenance",
		"Supply",
		"Consulting",
		"Construction",
		"Transportation",
	}

	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)
	var contracts []*entity.Contract

	// Generate 30 sample contracts
	for i := 0; i < 30; i++ {
		// Select a random contractor
		contractor := contractors[r.Intn(len(contractors))]

		// Generate contract dates
		contractDate := time.Now().AddDate(0, -r.Intn(6), -r.Intn(30))
		startDate := contractDate.AddDate(0, 0, r.Intn(30))
		endDate := startDate.AddDate(0, r.Intn(24), 0)

		// Generate a unique contract number
		contractNumber := fmt.Sprintf("CNT-%04d", i+1)

		contract := &entity.Contract{
			ContractNumber:  contractNumber,
			ContractDate:    contractDate,
			StartDate:       startDate,
			EndDate:         endDate,
			ContractAmount:  int64(r.Intn(1000000) * 1000), // Random amount up to 1 billion
			ContractType:    contractTypes[r.Intn(len(contractTypes))],
			ContractorID:    contractor.ID,
			OperationPeriod: r.Intn(36) + 1, // 1-36 months
			EquipmentPeriod: r.Intn(24) + 1, // 1-24 months
			Description:     faker.Paragraph(),
		}

		// Create the contract
		createdContract, err := cr.Create(context.Background(), contract)
		if err != nil {
			return nil, err
		}
		contracts = append(contracts, createdContract)
	}

	return contracts, nil
}
