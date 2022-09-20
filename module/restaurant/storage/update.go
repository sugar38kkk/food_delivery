package restaurantstorage

import (
	"context"
	"food-delivery/common"
	restaurantmodel "food-delivery/module/restaurant/model"
	"gorm.io/gorm"
)

func (store *sqlStore) UpdateDataWithConditions(
	context context.Context,
	updateData *restaurantmodel.RestaurantUpdate,
	conditions map[string]interface{},
) error {
	updateResult := store.db.Where(conditions).Updates(updateData)
	if updateResult.Error != nil {
		return updateResult.Error
	}
	return nil
}

func (store *sqlStore) IncreasedLikeCount(context context.Context, id int) error {
	db := store.db
	if err := db.Table(restaurantmodel.Restaurant{}.TableName()).Where("id = ?", id).
		Update("liked_count", gorm.Expr("liked_count + 1")).Error; err != nil {
		return common.ErrDB(err)
	}
	return nil
}

func (store *sqlStore) DecreasedLikeCount(context context.Context, int int) error {
	db := store.db
	if err := db.Table(restaurantmodel.Restaurant{}.TableName()).Where("id = ?", int).
		Update("liked_count", gorm.Expr("liked_count - 1")).Error; err != nil {
		return common.ErrDB(err)
	}
	return nil
}
