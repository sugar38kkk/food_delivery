package main

import (
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

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalln(err)
	}

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// POST /restaurants
	v1 := r.Group("/v1")

	restaurants := v1.Group("/restaurants")

	restaurants.POST("", func(c *gin.Context) {
		var data Restaurant

		if err := c.ShouldBind(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

			return
		}

		db.Create(&data)

		c.JSON(http.StatusOK, gin.H{
			"data": data,
		})
	})

	restaurants.GET("/", func(c *gin.Context) {

		var restaurants []Restaurant

		type Paging struct {
			Page  int `json:"page" form:"page"`
			Limit int `json:"limit" form:"limit"`
		}

		var pagingData Paging

		if err := c.ShouldBind(&pagingData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

			return
		}

		if pagingData.Page <= 0 {
			pagingData.Page = 1
		}

		if pagingData.Limit <= 0 {
			pagingData.Limit = 5
		}

		db.Offset((pagingData.Page - 1) * pagingData.Limit).
			Order("id desc").
			Limit(pagingData.Limit).
			Find(&restaurants)

		c.JSON(http.StatusOK, gin.H{
			"data": restaurants,
		})
	})

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

	restaurants.DELETE("/:id", func(c *gin.Context) {

		id, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

			return
		}

		db.Table(Restaurant{}.TableName()).Where("id = ?", id).Delete(nil)

		c.JSON(http.StatusOK, gin.H{
			"data": 1,
		})
	})

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
