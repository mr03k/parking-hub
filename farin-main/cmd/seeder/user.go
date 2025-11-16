package seeder

import (
	"context"
	"farin/domain/entity"
	"farin/domain/repository"
	"github.com/bxcodec/faker/v4"
	"math/rand"
	"time"
)

func SeedUsers(ur *repository.UserRepository) ([]*entity.User, error) {
	//roles := []entity.Role{entity.RoleAdmin, entity.RoleDriver}
	genders := []entity.Gender{entity.GenderMale, entity.GenderFemale}
	statuses := []entity.Status{entity.StatusInactive, entity.StatusActive}

	s := rand.NewSource(int64(time.Now().Nanosecond()))
	r := rand.New(s)
	var users []*entity.User

	for i := 0; i < 30; i++ {
		user := &entity.User{
			Username:     faker.Username(),
			Password:     "a1234567",
			FirstName:    faker.FirstName(),
			LastName:     faker.LastName(),
			Email:        faker.Email(),
			PhoneNumber:  faker.Phonenumber(),
			NationalID:   faker.Word(),
			PostalCode:   faker.Word(),
			CompanyName:  faker.ChineseFirstName(),
			ProfileImage: faker.URL(),
			Gender:       genders[r.Intn(len(genders)-1)],
			Address:      faker.Sentence(),
			Status:       statuses[r.Intn(len(statuses))],
			//Role:         roles[r.Intn(len(roles))],
		}
		user, err := ur.Create(context.Background(), user)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}
