package dal

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	mongoDBConnectionStringVar = "MONGODB_CONNECTION_STRING"
	mongoDBDatabaseVar         = "MONGODB_DATABASE"
	mongoDBCollectionVar       = "MONGODB_COLLECTION"
)

func ConnectDB() *mongo.Collection {
	mongoDBConnectionString := os.Getenv(mongoDBConnectionStringVar)
	if mongoDBConnectionString == "" {
		log.Fatal("MONGODB_CONNECTION_STRING not set")
	}

	mongoDBDatabase := os.Getenv(mongoDBDatabaseVar)
	if mongoDBDatabase == "" {
		log.Fatal("MONGODB_DATABASE not set")
	}

	mongoDBCollection := os.Getenv(mongoDBCollectionVar)
	if mongoDBCollection == "" {
		log.Fatal("MONGODB_COLLECTION not set")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(mongoDBConnectionString)

	c, err := mongo.NewClient(clientOptions)
	err = c.Connect(ctx)

	if err != nil {
		log.Fatal(err)
	}
	err = c.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	collection := c.Database(mongoDBDatabase).Collection(mongoDBCollection)

	return collection
}
