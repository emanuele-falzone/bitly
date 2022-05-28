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

type MongoRepository struct {
	client mongo.Client
}

type mongoEvent struct {
	Key      string `bson:"key"`
	Type     string `bson:"type"`
	DateTime string `bson:"datetime"`
}

func NewMongoRepository(connectionString string) (*MongoRepository, error) {
	// Create new mongo client with the given connection string
	client, err := mongo.NewClient(options.Client().ApplyURI(connectionString))
	if err != nil {
		return nil, err
	}

	// Connect with the mongo instance
	err = client.Connect(context.Background())
	if err != nil {
		return nil, err
	}

	// Return new EventRepository
	return &MongoRepository{client: *client}, nil
}

func (r *MongoRepository) Create(ctx context.Context, a Event) error {
	// Select database
	db := r.client.Database(db)

	// Select collection
	eventCollection := db.Collection(collection)

	// Insert event into collection
	_, err := eventCollection.InsertOne(context.Background(), mongoEvent{
		Key:      a.Redirection.Key,
		Type:     string(a.Type),
		DateTime: a.DateTime,
	})

	// Check for error during insert
	if err != nil {
		// There was some problem with mongo return error
		return &internal.Error{
			Code: internal.ErrInternal,
			Op:   "EventRepository: Create",
			Err:  err,
		}
	}

	// Return nil to signal that the operation was completed successfully
	return nil
}

func (r *MongoRepository) FindByRedirection(ctx context.Context, a redirection.Redirection) ([]Event, error) {
	// Select database
	db := r.client.Database(db)

	// Select collection
	eventCollection := db.Collection(collection)

	// Filter events by key
	filterCursor, err := eventCollection.Find(context.TODO(), bson.M{"key": a.Key})
	if err != nil {
		return nil, &internal.Error{
			Code: internal.ErrInternal,
			Op:   "EventRepository: FindByRedirection",
			Err:  err,
		}
	}

	// Materialize event slice
	var events []mongoEvent
	if err = filterCursor.All(context.TODO(), &events); err != nil {
		return nil, &internal.Error{
			Code: internal.ErrInternal,
			Op:   "EventRepository: FindByRedirection",
			Err:  err,
		}
	}

	// Map Event into Event
	results := make([]Event, len(events))

	for i, value := range events {
		results[i] = New(value.DateTime, Type(value.Type), a)
	}

	// Check result size
	if len(results) == 0 {
		return nil, &internal.Error{
			Code: internal.ErrNotFound,
			Op:   "EventRepository: FindByRedirection",
		}
	}

	// Return slice of domain objects
	return results, nil
}
