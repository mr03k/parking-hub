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

func SeedContractors(cr *repository.ContractorRepository) ([]*entity.Contractor, error) {
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
	var contractors []*entity.Contractor

	// Generate 20 sample contractors
	for i := 0; i < 20; i++ {
		// Generate a unique contractor code
		codeContractor := fmt.Sprintf("CON-%04d", i+1)

		contractor := &entity.Contractor{
			ContractorName:        faker.Word(),
			CodeContractor:        codeContractor,
			NumberRegistration:    faker.Word(),
			PersonContact:         faker.Name(),
			CEOName:               faker.Name(),
			SignatoriesAuthorized: faker.Sentence(),
			PhoneNumber:           faker.Phonenumber(),
			Email:                 faker.Email(),
			Address:               faker.Sentence(),
			TypeContract:          contractTypes[r.Intn(len(contractTypes))],
			NumberAccountBank:     faker.Word(), // This could be a more structured bank account number generator
			Description:           faker.Paragraph(),
		}

		// Create the contractor
		createdContractor, err := cr.Create(context.Background(), contractor)
		if err != nil {
			return nil, err
		}
		contractors = append(contractors, createdContractor)
	}

	return contractors, nil
}
