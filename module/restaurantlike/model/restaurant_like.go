package restaurantlikemodel

import (
	"fmt"
	"food-delivery/common"
	"time"
)

const EntityName = "UserLikeRestaurant"

type Like struct {
	RestaurantId int                `json:"restaurant_id" gorm:"restaurant_id;"`
	UserId       int                `json:"user_id" gorm:"user_id;"`
	CreatedAt    *time.Time         `json:"created_at" gorm:"created_at;"`
	User         *common.SimpleUser `json:"user" gorm:"foreignKey:UserId"`
}

func (Like) TableName() string {
	return "restaurant_likes"
}

func ErrCannotLikeRestaurant(err error) *common.AppError {
	return common.NewCustomError(
		err,
		fmt.Sprintf("Can not like this restaurant"),
		fmt.Sprintf("ErrCannotLikeRestaurant"),
	)
}

func ErrCannotDislikeRestaurant(err error) *common.AppError {
	return common.NewCustomError(
		err,
		fmt.Sprintf("Can not dislike this restaurant"),
		fmt.Sprintf("ErrCannotDislikeRestaurant"),
	)
}
