package restaurantlikestorage

import (
	"context"
	"fmt"
	"food-delivery/common"
	restaurantlikemodel "food-delivery/module/restaurantlike/model"
	"github.com/btcsuite/btcutil/base58"
	"time"
)

const timeLayout = "2006-01-02T15:04:05.999999"

func (store *sqlStore) GetRestaurantLikes(
	context context.Context,
	ids []int,
) (map[int]int, error) {
	db := store.db
	userLikes := make(map[int]int)
	type sqlData struct {
		RestaurantId int `gorm:"column:restaurant_id;"`
		LikeCount    int `gorm:"column:count;"`
	}
	var listLike []sqlData
	if err := db.Table(restaurantlikemodel.Like{}.TableName()).
		Select("restaurant_id, count(restaurant_id) as count").
		Where("restaurant_id in (?)", ids).
		Group("restaurant_id").
		Find(&listLike).Error; err != nil {
		return userLikes, common.ErrDB(err)
	}
	for _, item := range listLike {
		userLikes[item.RestaurantId] = item.LikeCount
	}
	return userLikes, nil
}

func (store *sqlStore) CheckUserHasLiked(context context.Context, userId int, resIds []int) (map[int]bool, error) {
	db := store.db
	userHasLiked := make(map[int]bool)
	type sqlData struct {
		RestaurantId int `gorm:"column:restaurant_id;"`
	}
	var result []sqlData
	if err := db.Table(restaurantlikemodel.Like{}.TableName()).
		Where("user_id = ? AND restaurant_id IN ?", userId, resIds).
		Find(&result).Error; err != nil {
		return userHasLiked, err
	}
	for _, item := range result {
		userHasLiked[item.RestaurantId] = true
	}
	return userHasLiked, nil
}

func (store *sqlStore) GetUsersLikeRestaurant(
	context context.Context,
	conditions map[string]interface{},
	filter *restaurantlikemodel.Filter,
	paging *common.Paging,
) ([]common.SimpleUser, error) {
	db := store.db
	var result []restaurantlikemodel.Like
	db = db.Table(restaurantlikemodel.Like{}.TableName()).Where(conditions)
	if v := filter; filter != nil {
		if v.RestaurantId > 0 {
			db = db.Where("restaurant_id = ?", v.RestaurantId)
		}
	}
	if paging != nil {
		if err := db.Count(&paging.Total).Error; err != nil {
			return nil, common.ErrDB(err)
		}
		db = db.Offset(paging.Offset()).Limit(paging.Limit)
	}

	db = db.Preload("User")

	if v := paging.FakeCursor; v != "" {
		timeCreated, err := time.Parse(timeLayout, string(base58.Decode(v)))

		if err != nil {
			return nil, common.ErrDB(err)
		}

		db = db.Where("created_at < ?", timeCreated.Format("2006-01-02 15:04:05"))
	} else {
		db = db.Offset(paging.Offset())
	}
	if err := db.Order("created_at desc").Find(&result).Error; err != nil {
		return nil, common.ErrDB(err)
	}

	listUsers := make([]common.SimpleUser, len(result))
	for i, item := range result {
		result[i].User.CreatedAt = item.CreatedAt
		result[i].User.UpdatedAt = nil
		listUsers[i] = *result[i].User

		if i == len(result)-1 {
			cursorStr := base58.Encode([]byte(fmt.Sprintf("%v", item.CreatedAt.Format(timeLayout))))
			paging.NextCursor = cursorStr
		}
	}
	return listUsers, nil
}
