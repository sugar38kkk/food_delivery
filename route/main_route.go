package route

import (
	"food-delivery/component/appctx"
	"food-delivery/middleware"
	restaurantmodel "food-delivery/module/restaurant/model"
	"food-delivery/module/restaurant/transport/ginrestaurant"
	"food-delivery/module/upload/transport/ginupload"
	"food-delivery/module/user/transport/ginuser"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func SetupRoute(appContext appctx.AppContext, v1 *gin.RouterGroup) {
	v1.POST("/upload", ginupload.UploadImage(appContext))

	restaurants := v1.Group("/restaurants", middleware.RequiredAuth(appContext))
	auth := v1.Group("/auth")

	//Auth
	auth.POST("/register", ginuser.Register(appContext))
	auth.POST("/login", ginuser.Login(appContext))
	auth.GET("/profile", middleware.RequiredAuth(appContext), ginuser.Profile(appContext))

	//Restaurant

	restaurants.POST("/", ginrestaurant.CreateRestaurant(appContext))

	restaurants.GET("/", ginrestaurant.ListRestaurant(appContext))

	restaurants.GET("/:id", func(c *gin.Context) {

		id, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

			return
		}

		var myRestaurant restaurantmodel.Restaurant

		appContext.GetMainDBConnection().Where("id = ?", id).First(&myRestaurant)

		c.JSON(http.StatusOK, gin.H{
			"data": myRestaurant,
		})
	})

	restaurants.PATCH("/:id", func(c *gin.Context) {

		id, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

			return
		}

		var data restaurantmodel.RestaurantUpdate

		if err := c.ShouldBind(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

			return
		}

		appContext.GetMainDBConnection().Where("id = ?", id).Updates(&data)

		c.JSON(http.StatusOK, gin.H{
			"data": data,
		})
	})

	restaurants.DELETE("/:id", ginrestaurant.DeleteRestaurant(appContext))
}
