package route

import (
	"food-delivery/component/appctx"
	"food-delivery/middleware"
	"food-delivery/module/user/transport/ginuser"
	"github.com/gin-gonic/gin"
)

func SetupAdminRoute(appContext appctx.AppContext, v1 *gin.RouterGroup) {
	admin := v1.Group("/admin", middleware.RequiredAuth(appContext), middleware.CheckRole(appContext, "admin", "mod"))

	{
		admin.GET("/profile", ginuser.Profile(appContext))
	}
}
