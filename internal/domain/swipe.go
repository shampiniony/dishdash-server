package domain

type Swipe struct {
	LobbyID string    `json:"lobbyId"`
	CardID  int64     `json:"cardId"`
	UserID  string    `json:"userId"`
	Type    SwipeType `json:"type"`
}

type SwipeType string

var (
	LIKE    SwipeType = "like"
	DISLIKE SwipeType = "dislike"
)
