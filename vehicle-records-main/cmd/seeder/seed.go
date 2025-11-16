package seeder

import (
	"context"
	"git.abanppc.com/farin-project/vehicle-records/infrastructure/godotenv"
	gormdb "git.abanppc.com/farin-project/vehicle-records/infrastructure/gorm"
	"log"
)

func Seed(fakeData bool) error {
	env := godotenv.NewEnv()
	env.Load()
	ctx := context.Background()
	gorm := gormdb.NewGORMDB(env)
	if err := gorm.Setup(ctx); err != nil {
		log.Fatalf("failed to setup gorm:%s", err)
	}

	return nil
}
