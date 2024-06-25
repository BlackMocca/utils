package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Client struct {
	db               *mongo.Client
	connectionString string
}

func NewConnection(ctx context.Context, connectionString string, registry *bsoncodec.Registry) (*Client, error) {
	client, err := mongo.Connect(ctx,
		options.Client().ApplyURI(connectionString).SetRetryReads(true).SetRetryWrites(true).SetRegistry(registry),
	)
	if err != nil {
		return nil, err
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	return &Client{
		db:               client,
		connectionString: connectionString,
	}, nil
}

func (c *Client) GetClient() *mongo.Client {
	return c.db
}

func (c *Client) SetDB(db *mongo.Client) {
	c.db = db
}

func (c *Client) Close(ctx context.Context) error {
	return c.db.Disconnect(ctx)
}
