package ginuser

import (
	"food-delivery/common"
	"food-delivery/component/appctx"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Profile(appCtx appctx.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		u := c.MustGet(common.CurrentUser).(common.Requester)
		//newPass := "kdsjkdfsjkdjfksdf"
		//type update struct {
		//	NewPass *string
		//}
		//
		//log.Println( update{ NewPass: &newPass})

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(u))
	}
}
