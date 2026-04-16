package repository

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const connectionString = "mongodb://%s:%s@%s/?directConnection=true"

var client *mongo.Client

type mongoEnvArgs struct {
	host     string
	user     string
	password string
	dbName   string
}

func InitConnection() *mongo.Database {
	if err := godotenv.Load(); err != nil {
		logrus.WithError(err).Warn("Error loading .env file")
	}

	db, err := connect()
	if err != nil {
		logrus.WithError(err).Fatal("Failed to connect to database")
	}
	defer close()

	return db
}

func loadEnvData() mongoEnvArgs {
	mongoEnvArgs := mongoEnvArgs{
		host:     os.Getenv("MONGO_HOST"),
		user:     os.Getenv("MONGO_USER"),
		password: os.Getenv("MONGO_PASSWORD"),
		dbName:   os.Getenv("MONGO_DB"),
	}

	if mongoEnvArgs.host == "" {
		mongoEnvArgs.host = "localhost:27017"
	}
	return mongoEnvArgs
}

func connect() (*mongo.Database, error) {
	mongoData := loadEnvData()

	if mongoData.user == "" || mongoData.password == "" || mongoData.dbName == "" {
		return nil, fmt.Errorf("missing required environment variables: MONGO_USER, MONGO_PASSWORD, MONGO_DB")
	}

	uri := fmt.Sprintf(connectionString, mongoData.user, mongoData.password, mongoData.host)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	opts := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	pingCtx, pingCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer pingCancel()

	if err := client.Ping(pingCtx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	return client.Database(mongoData.dbName), nil
}

func close() error {
	if client == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return client.Disconnect(ctx)
}
