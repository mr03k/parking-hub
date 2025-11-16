package seeder

import (
	"farin/domain/entity"
	"farin/domain/repository"
	"fmt"
	"github.com/bxcodec/faker/v4"
	"math/rand"
	"time"
)

import (
	"context"
)

func SeedRings(
	rr *repository.RingRepository,
) ([]*entity.Ring, error) {
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)
	var rings []*entity.Ring

	// Generate 20 sample rings
	for i := 0; i < 20; i++ {
		// Generate a unique ring code
		ringCode := fmt.Sprintf("RNG-%04d", i+1)

		// Generate a ring name
		ringName := fmt.Sprintf("%s %s Ring", faker.Word(), faker.Word())

		// Generate a mock geometry (this is a simplified example)
		// In a real-world scenario, you'd use a proper GIS library to generate valid geometries
		geom := fmt.Sprintf("LINESTRING(%f %f, %f %f)",
			r.Float64()*100, r.Float64()*100, // Start point
			r.Float64()*100, r.Float64()*100) // End point

		ring := &entity.Ring{
			RingCode: ringCode,
			Length:   r.Float64() * 10000, // Random length up to 10000
			RingName: ringName,
			Geom:     geom,
		}

		// Create the ring
		createdRing, err := rr.Create(context.Background(), ring)
		if err != nil {
			return nil, err
		}
		rings = append(rings, createdRing)
	}

	return rings, nil
}
