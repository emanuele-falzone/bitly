package mongo

import (
	"context"

	"github.com/emanuelefalzone/bitly/internal"
	"github.com/emanuelefalzone/bitly/internal/domain/event"
	"github.com/emanuelefalzone/bitly/internal/domain/redirection"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	DB         = "bitly"
	COLLECTION = "events"
)

type MongoEventRepository struct {
	client mongo.Client
}

type Event struct {
	Key      string `bson:"key"`
	Type     string `bson:"type"`
	DateTime string `bson:"datetime"`
}

func NewEventRepository(connectionString string) (*MongoEventRepository, error) {
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

	// Return new MongoEventRepository
	return &MongoEventRepository{client: *client}, nil
}

func (r *MongoEventRepository) Create(ctx context.Context, a event.Event) error {
	// Select database
	db := r.client.Database(DB)

	// Select collection
	eventCollection := db.Collection(COLLECTION)

	// Insert event into collection
	_, err := eventCollection.InsertOne(context.Background(), Event{
		Key:      a.Redirection.Key,
		Type:     string(a.Type),
		DateTime: a.DateTime,
	})

	// Check for error during insert
	if err != nil {
		// There was some problem with mongo return error
		return &internal.Error{Code: internal.ErrInternal, Op: "MongoEventRepository: Create", Err: err}
	}

	// Return nil to signal that the operation was completed successfully
	return nil
}

func (r *MongoEventRepository) FindByRedirection(ctx context.Context, a redirection.Redirection) ([]event.Event, error) {
	// Select database
	db := r.client.Database(DB)

	// Select collection
	eventCollection := db.Collection(COLLECTION)

	// Filter events by key
	filterCursor, err := eventCollection.Find(context.TODO(), bson.M{"key": a.Key})
	if err != nil {
		return nil, &internal.Error{Code: internal.ErrInternal, Op: "MongoEventRepository: FindByRedirection", Err: err}
	}

	// Materialize event slice
	var events []Event
	if err = filterCursor.All(context.TODO(), &events); err != nil {
		return nil, &internal.Error{Code: internal.ErrInternal, Op: "MongoEventRepository: FindByRedirection", Err: err}
	}

	// Map Event into event.Event
	var results []event.Event
	for _, value := range events {
		results = append(results, event.New(value.DateTime, event.Type(value.Type), a))
	}

	// Check result size
	if len(results) == 0 {
		return nil, &internal.Error{Code: internal.ErrNotFound, Op: "MongoEventRepository: FindByRedirection"}
	}

	// Return slice of domain objects
	return results, nil
}
