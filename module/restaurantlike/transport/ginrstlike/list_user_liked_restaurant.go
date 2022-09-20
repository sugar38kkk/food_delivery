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

func ListUserLikedRestaurant(appCtx appctx.AppContext) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		uid, err := common.FromBase58(ctx.Param("id"))
		if err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		filter := restaurantlikemodel.Filter{RestaurantId: int(uid.GetLocalID())}
		var paging common.Paging
		if err := ctx.ShouldBind(&paging); err != nil {
			panic(common.ErrInvalidRequest(err))
		}
		//if err := paging.Preprocess(); err != nil {
		//	panic(common.ErrInvalidRequest(err))
		//}
		paging.Fulfill()
		store := restaurantlikestorage.NewSQLStore(appCtx.GetMainDBConnection())
		biz := restaurantlikebusiness.NewListUsersLikeRestaurantBiz(store)

		listUsers, err := biz.ListUsersLikeRestaurant(
			ctx.Request.Context(),
			&filter,
			&paging,
		)
		if err != nil {
			panic(err)
		}
		for i := range listUsers {
			listUsers[i].Mask(false)
		}
		ctx.JSON(http.StatusOK, common.NewSuccessResponse(listUsers, paging, filter))
	}
}
