package event

import (
	"context"

	"github.com/emanuelefalzone/bitly/internal"
	"github.com/emanuelefalzone/bitly/internal/application/redirection"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	db         = "bitly"
	collection = "events"
)

// MongoRepository is an event repository that store values in mongo
type MongoRepository struct {
	client mongo.Client
}

type mongoEvent struct {
	Key      string `bson:"key"`
	Type     string `bson:"type"`
	DateTime string `bson:"datetime"`
}

// NewMongoRepository creates a new event repository that store values in mongo
func NewMongoRepository(connection string) (*MongoRepository, error) {
	// Create new mongo client with the given connection string
	client, err := mongo.NewClient(options.Client().ApplyURI(connection))
	if err != nil {
		return nil, err
	}

	// Connect with the mongo instance
	if err = client.Connect(context.Background()); err != nil {
		return nil, err
	}

	// Return new MongoRepository
	return &MongoRepository{client: *client}, nil
}

func (repo *MongoRepository) Create(ctx context.Context, value *Event) error {
	// Select database
	db := repo.client.Database(db)

	// Select collection
	eventCollection := db.Collection(collection)

	// Insert event into collection
	_, err := eventCollection.InsertOne(context.Background(), mongoEvent{
		Key:      value.Redirection.Key,
		Type:     string(value.Type),
		DateTime: value.DateTime,
	})

	// Check for error during insert
	if err != nil {
		// There was some problem with mongo return error
		return &internal.Error{
			Code: internal.ErrInternal,
			Op:   "MongoRepository: Create",
			Err:  err,
		}
	}

	// Return nil to signal that the operation was completed successfully
	return nil
}

func (repo *MongoRepository) FindByRedirection(ctx context.Context, value *redirection.Redirection) ([]*Event, error) {
	// Select database
	db := repo.client.Database(db)

	// Select collection
	eventCollection := db.Collection(collection)

	// Filter events by key
	filterCursor, err := eventCollection.Find(context.TODO(), bson.M{"key": value.Key})
	if err != nil {
		return nil, &internal.Error{
			Code: internal.ErrInternal,
			Op:   "MongoRepository: FindByRedirection",
			Err:  err,
		}
	}

	// Materialize event slice
	var events []mongoEvent
	if err = filterCursor.All(context.TODO(), &events); err != nil {
		return nil, &internal.Error{
			Code: internal.ErrInternal,
			Op:   "MongoRepository: FindByRedirection",
			Err:  err,
		}
	}

	// Map Event into Event
	results := make([]*Event, len(events))

	for i, event := range events {
		results[i] = New(event.DateTime, Type(event.Type), value)
	}

	// Check result size
	if len(results) == 0 {
		return nil, &internal.Error{
			Code: internal.ErrNotFound,
			Op:   "MongoRepository: FindByRedirection",
		}
	}

	// Return slice of domain objects
	return results, nil
}
