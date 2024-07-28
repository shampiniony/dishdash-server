package place

import (
	"dishdash.ru/internal/usecase"
	"github.com/gin-gonic/gin"
)

func SetupHandlers(r *gin.RouterGroup, cases usecase.Cases) {
	placeGroup := r.Group("places")
	placeGroup.POST("", SavePlace(cases.Place))
	placeGroup.GET("", GetAllPlaces(cases.Place))
	placeGroup.POST("tags", CreateTag(cases.Tag))
	placeGroup.GET("tags", GetAllTags(cases.Tag))
}
