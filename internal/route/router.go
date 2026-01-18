package route

import (
	"github.com/gin-gonic/gin"
	"github.com/symphire/arcrank/internal/handler"
)

func SetupRouter(
	player *handler.PlayerHandler,
	leaderboard *handler.LeaderboardHandler,
	search *handler.SearchHandler,
	health *handler.HealthHandler,
) *gin.Engine {
	r := gin.Default()

	r.GET("/health", health.Check)

	players := r.Group("/players")
	{
		players.POST("/", player.CreatePlayer)
		players.GET("/:id", player.GetPlayer)
		players.PATCH("/:id", player.UpdatePlayer)
	}

	r.GET("/leaderboard/top", leaderboard.GetTop)
	r.GET("/search", search.Query)

	return r
}
