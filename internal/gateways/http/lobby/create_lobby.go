package lobby

import (
	"net/http"

	"dishdash.ru/internal/usecase"
	"github.com/gin-gonic/gin"
)

// CreateLobby godoc
// @Summary Create a lobby
// @Description Create a new lobby in the database
// @Tags lobbies
// @Accept  json
// @Produce  json
// @Schemes http https
// @Param lobby body usecase.LobbyInput true "Lobby data"
// @Success 200 {object} lobbyOutput "Saved lobby"
// @Failure 400 "Bad Request"
// @Failure 500 "Internal Server Error"
// @Router /lobbies [post]
func CreateLobby(lobbyUseCase usecase.Lobby) gin.HandlerFunc {
	return func(c *gin.Context) {
		var lobbyInput usecase.LobbyInput
		err := c.BindJSON(&lobbyInput)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		lobby, err := lobbyUseCase.CreateLobby(c, lobbyInput)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, lobbyToOutput(lobby))
	}
}
