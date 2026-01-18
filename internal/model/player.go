package model

import (
	"time"
)

type Player struct {
	ID        string    `json:"id" bson:"_id"`
	Username  string    `json:"username" bson:"username"`
	XP        int       `json:"xp" bson:"xp"`
	Level     int       `json:"level" bson:"level"`
	Score     int       `json:"score" bson:"score"`
	Region    string    `json:"region" bson:"region"`
	Class     string    `json:"class" bson:"class"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

type CreatePlayerInput struct {
	Username string `json:"username" binding:"required"`
	Region   string `json:"region" binding:"required"`
	Class    string `json:"class" binding:"required"`
}

type UpdatePlayerInput struct {
	XP    *int `json:"xp"`
	Level *int `json:"level"`
	Score *int `json:"score"`
}
