package domain

type Card struct {
	ID               int64      `json:"id"`
	Title            string     `json:"title"`
	ShortDescription string     `json:"shortDescription"`
	Description      string     `json:"description"`
	Image            string     `json:"image"`
	Location         Coordinate `json:"location"`
	Address          string     `json:"address"`
	PriceMin         int        `json:"priceMin"`
	PriceMax         int        `json:"priceMax"`
	Tags             []*Tag     `json:"tags"`
}
