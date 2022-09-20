package ginrstlike

import (
	"food-delivery/common"
	"food-delivery/component/appctx"
	restaurantlikebusiness "food-delivery/module/restaurantlike/business"
	restaurantlikemodel "food-delivery/module/restaurantlike/model"
	restaurantlikestorage "food-delivery/module/restaurantlike/storage"
	"github.com/gin-gonic/gin"
	"net/http"
)

func UserDislikeRestaurant(appCtx appctx.AppContext) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		uid, err := common.FromBase58(ctx.Param("id"))
		if err != nil {
			panic(common.ErrInvalidRequest(err))
		}
		requester := ctx.MustGet(common.CurrentUser).(common.Requester)
		data := restaurantlikemodel.Like{
			RestaurantId: int(uid.GetLocalID()),
			UserId:       requester.GetUserId(),
		}
		store := restaurantlikestorage.NewSQLStore(appCtx.GetMainDBConnection())
		//decStore := restaurantstorage.NewSQLStore(appCtx.GetMainDBConnection())
		biz := restaurantlikebusiness.NewUserDislikeRestaurantBiz(store, appCtx.GetPubSub())
		if err := biz.DislikeRestaurant(ctx.Request.Context(), data.UserId, data.RestaurantId); err != nil {
			panic(err)
		}
		ctx.JSON(http.StatusOK, common.SimpleSuccessResponse(true))
	}
}
