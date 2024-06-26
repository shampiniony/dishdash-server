package domain

import (
	"time"
)

type Lobby struct {
	ID            string         `json:"id"`
	CreatedAt     time.Time      `json:"createdAt"`
	Location      Coordinate     `json:"location"`
	LobbySettings *LobbySettings `json:"lobbySettings"`
	Cards         []*Card        `json:"cards"`
	Matches       []*Match       `json:"matches"`
	FinalVotes    []*FinalVote   `json:"finalVotes"`
	Swipes        []*Swipe       `json:"swipes"`
}
