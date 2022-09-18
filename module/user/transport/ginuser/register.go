package ginuser

import (
	"food-delivery/common"
	"food-delivery/component/appctx"
	usermodel "food-delivery/module/user/model"
	"net/http"

	userbiz "food-delivery/module/user/business"
	userstorage "food-delivery/module/user/storage"

	"food-delivery/component/hasher"

	"github.com/gin-gonic/gin"
)

func Register(appctx appctx.AppContext) func(*gin.Context) {
	return func(c *gin.Context) {

		var data usermodel.UserCreate

		if err := c.ShouldBind(&data); err != nil {
			panic(err)
		}

		db := appctx.GetMainDBConnection()

		store := userstorage.NewSQLStore(db)
		md5 := hasher.NewMd5Hash()

		biz := userbiz.NewRegisterBiz(store, md5)

		if err := biz.RegisterUser(c, &data); err != nil {
			panic(err)
		}

		data.Mask(false)

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(data.FakeId.String()))

	}
}
