package domain

type Match struct {
	ID      int64  `json:"id"`
	LobbyID string `json:"lobbyId"`
	CardID  int64  `json:"cardId"`
}
