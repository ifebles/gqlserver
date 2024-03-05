package database

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func Close(client *mongo.Client, ctx context.Context, cancel context.CancelFunc) {
	// CancelFunc to cancel to context
	defer cancel()

	// client provides a method to close a mongoDB connection.
	defer func() {
		// client.Disconnect method also has deadline.
		if err := client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
}

func Connect(uri string) (*mongo.Client, context.Context, context.CancelFunc, error) {
	// ctx will be used to set deadline for process, here deadline will of 30 seconds.
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	// mongo.Connect return mongo.Client method
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))

	// Confirmation of connection
	return client, ctx, cancel, err
}

// This method used to ping the mongoDB, return error if any.
func ping(client *mongo.Client, ctx context.Context) error {
	// mongo.Client has Ping to ping mongoDB, deadline of the Ping method will be determined by cxt
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return err
	}

	fmt.Println("Connected successfully")
	return nil
}

func GetConnection() {
	// Get Client, Context, CalcelFunc and err from connect method.
	client, ctx, cancel, err := Connect(os.Getenv("DB_URI"))

	if err != nil {
		panic(err)
	}

	defer Close(client, ctx, cancel)

	if err := ping(client, ctx); err != nil {
		panic(err)
	}
}

func SaveOne(client *mongo.Client, ctx context.Context, database, col string, doc any) (*mongo.InsertOneResult, error) {
	collection := client.Database(database).Collection(col)

	return collection.InsertOne(ctx, doc)
}

func Query(client *mongo.Client, ctx context.Context, database, col string, filter bson.D) (result *mongo.Cursor, err error) {
	collection := client.Database(database).Collection(col)

	return collection.Find(ctx, filter)
}
