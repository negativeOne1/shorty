package mongodb

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Client struct {
	*mongo.Client
	*mongo.Database
}

func New(ctx context.Context, conn string, dbName string) (*Client, error) {
	if conn == "" {
		return nil, fmt.Errorf("%v unset or empty", "MONGO_CONNECTION_STRING")
	}

	client, err := mongo.NewClient(options.Client().ApplyURI(conn))
	if err != nil {
		return nil, fmt.Errorf("failed to create a mongodb client: %v", err)
	}

	if err := client.Connect(ctx); err != nil {
		return nil, fmt.Errorf("failed to connect to mongo: %v", err)
	}

	return &Client{
		Client:   client,
		Database: client.Database(dbName),
	}, err
}

func (c *Client) CreateIndex(ctx context.Context, coll string, index mongo.IndexModel) error {
	db := c.Database.Collection(coll)

	_, err := db.Indexes().CreateOne(ctx, index)
	return err
}
