package restaurantlikebusiness

import (
	"context"
	"food-delivery/common"
	restaurantlikemodel "food-delivery/module/restaurantlike/model"
	"food-delivery/pubsub"
	"log"
)

type UserLikeRestaurantStore interface {
	//FindDataWithConditions(
	//	context context.Context,
	//	condition map[string]interface{},
	//) (*restaurantlikemodel.Like, error)
	Create(context context.Context, newData *restaurantlikemodel.Like) error
}

//type IncreaseLikeCountRstStore interface {
//	IncreasedLikeCount(context context.Context, id int) error
//}

type userLikeRestaurantBiz struct {
	store UserLikeRestaurantStore
	//incStore IncreaseLikeCountRstStore
	ps pubsub.Pubsub
}

func NewUserLikeRestaurantBiz(store UserLikeRestaurantStore, ps pubsub.Pubsub) *userLikeRestaurantBiz {
	return &userLikeRestaurantBiz{
		store: store,
		//incStore: incStore,
		ps: ps,
	}
}

func (biz *userLikeRestaurantBiz) LikeRestaurant(
	ctx context.Context,
	data *restaurantlikemodel.Like,
) error {
	err := biz.store.Create(ctx, data)
	if err != nil {
		return restaurantlikemodel.ErrCannotLikeRestaurant(err)
	}

	//send message
	if err := biz.ps.Publish(ctx, common.TopicUserLikeRestaurant, pubsub.NewMessage(data)); err != nil {
		log.Println(err)
	}

	//j := asyncjob.NewJob(func(ctx context.Context) error {
	//	return biz.incStore.IncreasedLikeCount(ctx, data.RestaurantId)
	//})
	//
	//asyncjob.NewGroup(true, j).Run(ctx)

	return nil
}
