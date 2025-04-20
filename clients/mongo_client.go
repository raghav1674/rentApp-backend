package clients

import (
	"context"
	"sample-web/configs"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MongoClient struct {
	Client   *mongo.Client
	Database *mongo.Database
}

func NewMongoClient(cfg configs.MongoConfig) (*MongoClient, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.TimeoutInSeconds)*time.Second)
	defer cancel()

	clientOpts := options.Client().ApplyURI(cfg.URI)
	clientOpts.Auth = &options.Credential{
		Username:   cfg.Username,
		Password:   cfg.Password,
		AuthSource: cfg.AuthSource,
	}
	client, err := mongo.Connect(clientOpts)
	if err != nil {
		return nil, err
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	return &MongoClient{
		Client:   client,
		Database: client.Database(cfg.Database),
	}, nil
}
