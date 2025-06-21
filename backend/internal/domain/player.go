package domain

import "time"

type Player struct {
	ID       string    `json:"id"`
	Name     string    `json:"name"`
	IsReady  bool      `json:"isReady"`
	JoinedAt time.Time `json:"joinedAt"`
	Score    int       `json:"score"`
}
