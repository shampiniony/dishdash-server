package domain

type LobbySettings struct {
	PriceMin    int     `json:"priceMin"`
	PriceMax    int     `json:"priceMax"`
	MaxDistance float64 `json:"maxDistance"`
	Tags        []Tag   `json:"tags"`
}
