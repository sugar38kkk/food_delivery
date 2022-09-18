package ginuser

import (
	"food-delivery/common"
	"food-delivery/component/appctx"
	"food-delivery/component/hasher"
	"food-delivery/component/tokenprovider/jwt"
	userbiz "food-delivery/module/user/business"
	usermodel "food-delivery/module/user/model"
	userstorage "food-delivery/module/user/storage"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Login(appCtx appctx.AppContext) func(*gin.Context) {
	return func(c *gin.Context) {
		var data usermodel.UserLogin

		if err := c.ShouldBind(&data); err != nil {
			panic(err)
		}

		db := appCtx.GetMainDBConnection()
		tokenProvider := jwt.NewTokenJWTProvider(appCtx.SecretKey())

		store := userstorage.NewSQLStore(db)
		md5 := hasher.NewMd5Hash()

		business := userbiz.NewLoginBiz(store, md5, tokenProvider, 60*60*24*30)
		account, err := business.Login(c.Request.Context(), &data)

		if err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(account))
	}
}
