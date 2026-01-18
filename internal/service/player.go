package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/symphire/arcrank/internal/db"
	"github.com/symphire/arcrank/internal/log"
	"github.com/symphire/arcrank/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type PlayerService struct {
	mongo *db.MongoDB
	es    *db.ElasticDB
	log   *log.Logger
}

func NewPlayerService(m *db.MongoDB, e *db.ElasticDB, l *log.Logger) *PlayerService {
	return &PlayerService{m, e, l}
}

// CreatePlayer creates the player in Mongo and immediately indexes into ES.
func (s *PlayerService) CreatePlayer(ctx context.Context, input model.CreatePlayerInput) (*model.Player, error) {
	now := time.Now().UTC()

	player := &model.Player{
		ID:        uuid.NewString(),
		Username:  input.Username,
		Region:    input.Region,
		Class:     input.Class,
		XP:        0,
		Level:     1,
		Score:     0,
		UpdatedAt: now,
	}

	if _, err := s.mongo.PlayerCollection.InsertOne(ctx, player); err != nil {
		s.log.Error("Mongo insert error", "error", err)
		return nil, err
	}

	if err := s.indexPlayer(ctx, player); err != nil {
		// TODO: rollback is needed to maintain consistency
		s.log.Error("Failed to index player in ES", "error", err)
		return nil, err
	}

	return player, nil
}

// indexPlayer sends the full player document to ES using the Mongo _id as ES document id.
func (s *PlayerService) indexPlayer(ctx context.Context, player *model.Player) error {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(player); err != nil {
		return fmt.Errorf("json encede error (player): %w", err)
	}

	res, err := s.es.Client.Index(
		s.es.Index,
		&buf,
		s.es.Client.Index.WithDocumentID(player.ID),
		s.es.Client.Index.WithContext(ctx),
	)
	if err != nil {
		return fmt.Errorf("ES index error (player): %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("ES index status error (player): %s", res.Status())
	}

	return nil
}

func (s *PlayerService) GetPlayer(ctx context.Context, id string) (*model.Player, error) {
	var player model.Player
	err := s.mongo.PlayerCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&player)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	}
	if err != nil {
		s.log.Error("Mongo find player error (get)", "error", err)
		return nil, err
	}
	return &player, nil
}

func (s *PlayerService) UpdatePlayer(ctx context.Context, id string, input model.UpdatePlayerInput) (*model.Player, error) {
	var player model.Player
	if err := s.mongo.PlayerCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&player); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		s.log.Error("Mongo find player error (update)", "error", err)
		return nil, err
	}

	// region Apply updates

	if input.XP != nil {
		player.XP = *input.XP
	}
	if input.Level != nil {
		player.Level = *input.Level
	}
	if input.Score != nil {
		player.Score = *input.Score
	}
	player.UpdatedAt = time.Now().UTC()

	// endregion

	if _, err := s.mongo.PlayerCollection.ReplaceOne(ctx, bson.M{"_id": id}, player); err != nil {
		s.log.Error("Mongo replace player error (update)", "error", err)
		return nil, err
	}

	if err := s.indexPlayer(ctx, &player); err != nil {
		s.log.Error("Failed to re-index player in ES", "error", err, "id", id)
	}

	return &player, nil
}
