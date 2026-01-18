package db

import (
	"context"
	"time"

	"github.com/symphire/arcrank/internal/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	Client           *mongo.Client
	PlayerCollection *mongo.Collection
}

func NewMongoDB(ctx context.Context, uri string, logger *log.Logger) (*MongoDB, error) {
	opts := options.Client().ApplyURI(uri)

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		logger.Error("Failed to connect to MongoDB", "error", err)
		return nil, err
	}

	if err = client.Ping(ctx, nil); err != nil {
		logger.Error("Failed to ping MongoDB", "error", err)
		return nil, err
	}

	// TODO: move hardcoded DB and collection to config file
	db := client.Database("arcrank")
	players := db.Collection("players")

	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "username", Value: 1}},
		Options: options.Index().SetUnique(true),
	}
	_, err = players.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		logger.Error("Failed to create player username index", "error", err)
		return nil, err
	}

	logger.Info("Connected to MongoDB with indexes ready", "db", "arcrank", "collection", "players")

	return &MongoDB{
		Client:           client,
		PlayerCollection: players,
	}, nil
}
