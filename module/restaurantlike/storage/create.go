package restaurantlikestorage

import (
	"context"
	"food-delivery/common"
	restaurantlikemodel "food-delivery/module/restaurantlike/model"
)

func (store *sqlStore) Create(context context.Context, newData *restaurantlikemodel.Like) error {
	db := store.db
	if err := db.Create(newData).Error; err != nil {
		return common.ErrDB(err)
	}
	return nil
}
