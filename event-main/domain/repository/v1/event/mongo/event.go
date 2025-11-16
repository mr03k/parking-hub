package mongo

import (
	"context"
	v1 "git.abanppc.com/farin-project/event/domain/entity/v1"
	"git.abanppc.com/farin-project/event/pkg/mongo"
	"github.com/google/uuid"
	mongoPKG "go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/acme"
)

type EventRepository struct {
	m *mongo.Mongo
}

// NewEventRepository creates a new instance of EventRepository for MongoDB
func NewEventRepository(m *mongo.Mongo) *EventRepository {
	return &EventRepository{
		m: m,
	}
}

// CreateUser inserts a new user into MongoDB
func (r *EventRepository) CreateEvent(ctx context.Context, event *v1.Event) error {
	event.ID = uuid.NewString()
	_, err := r.m.CE.InsertOne(ctx, event)
	if mongoPKG.IsDuplicateKeyError(err) {
		return acme.ErrAccountAlreadyExists
	}
	return err
}
