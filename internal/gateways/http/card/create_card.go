package card

import (
	"net/http"

	"dishdash.ru/internal/usecase"
	"github.com/gin-gonic/gin"
)

// CreateCard godoc
// @Summary Create a card
// @Description Create a new card in the database
// @Tags cards
// @Accept  json
// @Produce  json
// @Schemes http https
// @Param card body cardInput true "Card data"
// @Success 200 {object} cardOutput "Saved card"
// @Failure 400 "Bad Request"
// @Failure 500 "Internal Server Error"
// @Router /cards [post]
func CreateCard(cardUseCase *usecase.Card) gin.HandlerFunc {
	return func(c *gin.Context) {
		var cardInput cardInput
		err := c.BindJSON(&cardInput)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}

		card := inputToCard(cardInput)
		err = cardUseCase.CreateCard(c, card)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}
		err = cardUseCase.AttachTagsToCard(c, cardInput.Tags, card.ID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}
		card, err = cardUseCase.GetCardByID(c, card.ID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		}
		c.JSON(http.StatusOK, cardToOutput(card))
	}
}