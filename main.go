package main

import (
	"food-delivery/component/appctx"
	"food-delivery/component/uploadprovider"
	"food-delivery/middleware"
	"food-delivery/module/restaurant/transport/ginrestaurant"
	"food-delivery/module/upload/transport/ginupload"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
	"strconv"
)

const TableName string = "restaurants"

type Restaurant struct {
	Id   int    `json:"id" gorm:"column:id;"`
	Name string `json:"name" gorm:"column:name"`
	Addr string `json:"addr" gorm:"column:addr"`
}

func (Restaurant) TableName() string { return TableName }

type RestaurantUpdate struct {
	Name *string `json:"name" gorm:"column:name"`
	Addr *string `json:"addr" gorm:"column:addr"`
}

func (RestaurantUpdate) TableName() string { return TableName }

func main() {
	dsn := os.Getenv("MYSQL_CONN_STRING")

	s3BucketName := os.Getenv("S3BucketName")
	s3Region := os.Getenv("S3Region")
	s3APIKey := os.Getenv("S3APIKey")
	s3SecretKey := os.Getenv("S3SecretKey")
	s3Domain := os.Getenv("S3Domain")
	//secretKey := os.Getenv("SYSTEM_SECRET")

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalln(err)
	}

	db = db.Debug()

	s3Provider := uploadprovider.NewS3Provider(s3BucketName, s3Region, s3APIKey, s3SecretKey, s3Domain)

	appContext := appctx.NewAppContext(db, s3Provider)

	r := gin.Default()
	r.Use(middleware.Recover(appContext))
	r.Static("/static", "static")

	// POST /restaurants
	v1 := r.Group("/v1")

	v1.POST("/upload", ginupload.UploadImage(appContext))

	restaurants := v1.Group("/restaurants")

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

		var myRestaurant Restaurant

		db.Where("id = ?", id).First(&myRestaurant)

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

		var data RestaurantUpdate

		if err := c.ShouldBind(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

			return
		}

		db.Where("id = ?", id).Updates(&data)

		c.JSON(http.StatusOK, gin.H{
			"data": data,
		})
	})

	restaurants.DELETE("/:id", ginrestaurant.DeleteRestaurant(appContext))

	r.Run()

	//newRestaurant := Restaurant{Name: "Big Mouth", Addr: "K09/54 Ha Van Tri"}
	//
	//if err := db.Create(&newRestaurant).Error; err != nil {
	//	log.Println(err)
	//}

	//var myRestaurant Restaurant
	//
	//if err := db.Where("id = ?", 3).First(&myRestaurant).Error; err != nil {
	//	log.Println(err)
	//}
	//
	//log.Println(myRestaurant)
	//
	//newName := "Sugar Tech"
	//updateData := RestaurantUpdate{Name: &newName}
	//
	//if err := db.Where("id = ?", 3).Updates(&updateData).Error; err != nil {
	//	log.Println(err)
	//}
	//
	//log.Println(myRestaurant)
	//
	//if err := db.Table(Restaurant{}.TableName()).Where("id = ?", 2).Delete(nil).Error; err != nil {
	//	log.Println(err)
	//}

}
