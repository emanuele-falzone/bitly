package util

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ClearMongo(ctx context.Context, connectionString, database string) error {
	// Create new mongo client with the given connection string
	client, err := mongo.NewClient(options.Client().ApplyURI(connectionString))
	if err != nil {
		return err
	}

	// Connect with the mongo instance
	err = client.Connect(context.Background())
	if err != nil {
		return err
	}

	// Select database
	db := client.Database(database)

	// Drop db
	return db.Drop(context.Background())
}
