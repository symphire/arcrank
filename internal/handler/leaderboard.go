package handler

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/symphire/arcrank/internal/service"
)

type LeaderboardHandler struct {
	service *service.LeaderboardService
}

func NewLeaderboardHandler(service *service.LeaderboardService) *LeaderboardHandler {
	return &LeaderboardHandler{service: service}
}

func (h *LeaderboardHandler) GetTop(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "100")
	limit, _ := strconv.Atoi(limitStr)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	players, err := h.service.GetTop(ctx, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "leaderboard query failed", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, players)
}
