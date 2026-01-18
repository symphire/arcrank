package main

import (
	"context"
	"os"

	"github.com/symphire/arcrank/config"
	"github.com/symphire/arcrank/internal/db"
	"github.com/symphire/arcrank/internal/handler"
	"github.com/symphire/arcrank/internal/log"
	"github.com/symphire/arcrank/internal/route"
	"github.com/symphire/arcrank/internal/service"
)

func main() {
	ctx := context.Background()

	// ---- Load config ----
	cfg := config.Load()

	// ---- Init logger
	logger := log.New(cfg.LogLevel)
	logger.Info("Starting ArcRank server...")

	// ---- Connect DB ----
	mongo, err := db.NewMongoDB(ctx, cfg.MongoURL, logger)
	if err != nil {
		logger.Fatal("Failed to connect to MongoDB", "error", err)
	}
	es, err := db.NewElasticDB(cfg.ElasticURL, cfg.ElasticIndex, logger)
	if err != nil {
		logger.Fatal("Failed to connect to Elasticsearch", "error", err)
	}

	// ---- Services ----
	playerService := service.NewPlayerService(mongo, es, logger)
	leaderboardService := service.NewLeaderboardService(es, logger)
	searchService := service.NewSearchService(es, logger)

	// ---- Handlers ----
	playerHandler := handler.NewPlayerHandler(playerService)
	leaderboardHandler := handler.NewLeaderboardHandler(leaderboardService)
	searchHandler := handler.NewSearchHandler(searchService)
	healthHandler := handler.NewHealthHandler()

	// ---- Router ----
	r := route.SetupRouter(
		playerHandler,
		leaderboardHandler,
		searchHandler,
		healthHandler,
	)

	// ---- Run ----
	if err := r.Run(cfg.RunPort); err != nil {
		logger.Error("Server exited with error", "error", err)
		os.Exit(1)
	}
}
