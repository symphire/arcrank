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

type SearchService struct {
	es  *db.ElasticDB
	log *log.Logger
}

func NewSearchService(es *db.ElasticDB, l *log.Logger) *SearchService {
	return &SearchService{es: es, log: l}
}

func (s *SearchService) SearchByUsername(ctx context.Context, q string, limit int) ([]model.Player, error) {
	if q == "" {
		return []model.Player{}, nil
	}
	if limit <= 0 {
		limit = 50
	}

	query := map[string]any{
		"size": limit,
		"query": map[string]any{
			"match_phrase_prefix": map[string]any{
				"username": q,
			},
		},
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, fmt.Errorf("encode search query: %w", err)
	}

	res, err := s.es.Client.Search(
		s.es.Client.Search.WithIndex(s.es.Index),
		s.es.Client.Search.WithBody(&buf),
		s.es.Client.Search.WithContext(ctx),
	)
	if err != nil {
		return nil, fmt.Errorf("ES search error (by username): %w", err)
	}
	defer res.Body.Close()

	var esRes struct {
		Hits struct {
			Hits []struct {
				Source model.Player `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	if err := json.NewDecoder(res.Body).Decode(&esRes); err != nil {
		return nil, fmt.Errorf("decode search response: %w", err)
	}

	out := make([]model.Player, 0, len(esRes.Hits.Hits))
	for _, player := range esRes.Hits.Hits {
		out = append(out, player.Source)
	}

	return out, nil
}
