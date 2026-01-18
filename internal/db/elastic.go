package db

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/symphire/arcrank/internal/log"
)

type ElasticDB struct {
	Client *elasticsearch.Client
	Index  string
}

func NewElasticDB(uri, index string, logger *log.Logger) (*ElasticDB, error) {
	cfg := elasticsearch.Config{
		Addresses: []string{uri},
		Transport: &http.Transport{
			Proxy: nil,
		},
	}

	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		logger.Error("Failed to create Elasticsearch client", "error", err)
		return nil, err
	}

	res, err := client.Info()
	if err != nil {
		logger.Error("Failed to connect to Elasticsearch", "error", err)
		return nil, err
	}
	defer res.Body.Close()

	logger.Info("Connected to Elasticsearch")

	es := &ElasticDB{
		Client: client,
		Index:  index,
	}

	if err := es.ensureIndex(logger); err != nil {
		return nil, err
	}

	return es, nil
}

func (e *ElasticDB) ensureIndex(logger *log.Logger) error {
	ctx := context.Background()

	existRes, err := e.Client.Indices.Exists([]string{e.Index}, e.Client.Indices.Exists.WithContext(ctx))
	if err != nil {
		logger.Error("Failed to check ES index existence", "error", err)
		return err
	}
	defer existRes.Body.Close()

	if existRes.StatusCode == http.StatusOK {
		logger.Info("Elasticsearch index exists", "index", e.Index)
		return nil
	}

	mapping := map[string]any{
		"settings": map[string]any{
			"number_of_shards":   1,
			"number_of_replicas": 0,
		},
		"mappings": map[string]any{
			"properties": map[string]any{
				"username":   map[string]any{"type": "text"},
				"score":      map[string]any{"type": "integer"},
				"level":      map[string]any{"type": "integer"},
				"region":     map[string]any{"type": "keyword"},
				"class":      map[string]any{"type": "keyword"},
				"updated_at": map[string]any{"type": "date"},
			},
		},
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(mapping); err != nil {
		logger.Error("Failed to encode ES mapping", "error", err)
		return err
	}

	createRes, err := e.Client.Indices.Create(
		e.Index,
		e.Client.Indices.Create.WithBody(&buf),
		e.Client.Indices.Create.WithContext(ctx),
	)
	if err != nil {
		logger.Error("Failed to create ES index", "error", err)
		return err
	}
	defer createRes.Body.Close()

	body, _ := io.ReadAll(createRes.Body)
	if createRes.IsError() {
		err := fmt.Errorf("failed to create ES index: status=%s body=%s", createRes.Status(), string(body))
		logger.Error("Failed to create ES index", "error", err)
		return err
	}

	logger.Info("Created ES index", "index", e.Index)

	return nil
}
