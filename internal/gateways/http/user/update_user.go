package user

import (
	"net/http"

	"dishdash.ru/internal/domain"

	"dishdash.ru/internal/usecase"
	"github.com/gin-gonic/gin"
)

// UpdateUser godoc
// @Summary Update a user
// @Description Update an existing user in the database
// @Tags users
// @Accept  json
// @Produce  json
// @Schemes http https
// @Param user body domain.User true "User data"
// @Success 200 {object} domain.User "Updated user"
// @Failure 400 "Bad Request"
// @Failure 500 "Internal Server Error"
// @Router /users [put]
func UpdateUser(userUseCase usecase.User) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := new(domain.User)
		err := c.BindJSON(&user)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		user, err = userUseCase.UpdateUser(c, user)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, user)
	}
}
