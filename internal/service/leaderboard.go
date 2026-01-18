package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/symphire/arcrank/internal/db"
	"github.com/symphire/arcrank/internal/log"
	"github.com/symphire/arcrank/internal/model"
)

type LeaderboardService struct {
	es  *db.ElasticDB
	log *log.Logger
}

func NewLeaderboardService(es *db.ElasticDB, l *log.Logger) *LeaderboardService {
	return &LeaderboardService{es: es, log: l}
}

func (s *LeaderboardService) GetTop(ctx context.Context, limit int) ([]model.Player, error) {
	if limit <= 0 {
		limit = 100
	}

	query := map[string]any{
		"size": limit,
		"sort": []map[string]any{
			{"score": map[string]any{"order": "desc"}},
		},
		"query": map[string]any{
			"match_all": map[string]any{},
		},
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, fmt.Errorf("encode leaderboard query: %w", err)
	}

	res, err := s.es.Client.Search(
		s.es.Client.Search.WithIndex(s.es.Index),
		s.es.Client.Search.WithBody(&buf),
		s.es.Client.Search.WithContext(ctx),
	)
	if err != nil {
		return nil, fmt.Errorf("search leaderboard: %w", err)
	}
	defer res.Body.Close()

	var esResp struct {
		Hits struct {
			Hits []struct {
				Source model.Player `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	if err := json.NewDecoder(res.Body).Decode(&esResp); err != nil {
		return nil, fmt.Errorf("decode leaderboard response: %w", err)
	}

	players := make([]model.Player, 0, len(esResp.Hits.Hits))
	for _, h := range esResp.Hits.Hits {
		players = append(players, h.Source)
	}

	return players, nil
}
