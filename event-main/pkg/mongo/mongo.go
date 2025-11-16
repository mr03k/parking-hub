package mongo

import (
	"context"
	"errors"
	"git.abanppc.com/farin-project/event/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type Mongo struct {
	CE          *mongo.Collection //event collection
	DB          *mongo.Database
	mongoConfig *config.MongoDB
}

func NewMongo(mongoConfig *config.MongoDB) *Mongo {
	return &Mongo{
		mongoConfig: mongoConfig,
	}
}

func (m *Mongo) Setup(ctx context.Context) error {
	if m.DB != nil {
		m.DB.Client().Disconnect(ctx)
	}
	// Set up MongoDB client options
	clientOptions := options.Client().ApplyURI(m.mongoConfig.Host)

	// Connect to MongoDB
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return err
	}
	db := client.Database(m.mongoConfig.DBName)
	m.DB = db

	cEvent := db.Collection("event")
	indexModels := []mongo.IndexModel{
		{
			Keys: bson.M{
				"lprVehicleRTKLocation": "2dsphere",
			},
		},
		{
			Keys: bson.M{
				"lprVehicleGPSLocation": "2dsphere",
			},
		},
	}
	_, err = cEvent.Indexes().CreateMany(ctx, indexModels)
	if err != nil {
		return err
	}
	m.CE = cEvent
	return nil
}

// HealthCheck checks if MongoDB is ready by performing a basic ping and a small query
func (m *Mongo) HealthCheck() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	// Check if the collection is initialized
	if m.CE == nil {
		return errors.New("mongo collection is not initialized")
	}

	// Perform a simple operation to ensure MongoDB is ready
	// Here we can run a simple find query that checks if MongoDB is responsive
	if err := m.CE.Database().Client().Ping(ctx, nil); err != nil {
		return errors.New("failed to ping MongoDB: " + err.Error())
	}

	// Optionally, you could perform a lightweight query (e.g., count documents in the collection)
	_, err := m.CE.CountDocuments(ctx, bson.M{})
	if err != nil {
		return errors.New("mongoDB query failed: " + err.Error())
	}

	return nil
}
