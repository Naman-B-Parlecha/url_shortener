package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Connect() (*mongo.Client, error) {
	godotenv.Load()
	val, exist := os.LookupEnv("MONGO_URI")
	if !exist {
		log.Fatal("MONGO_URI environment variable not set")
		return nil, fmt.Errorf("MONGO_URI environment variable not set")
	}

	uri := val

	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
		return nil, fmt.Errorf("failed to create MongoDB client: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Println("Failed to connect to MongoDB")
		fmt.Println(err)
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	if err == nil {
		fmt.Println("Successfully connected to MongoDB")
	}

	return client, nil
}
