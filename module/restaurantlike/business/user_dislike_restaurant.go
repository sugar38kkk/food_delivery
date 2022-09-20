package restaurantlikebusiness

import (
	"context"
	"food-delivery/common"
	restaurantlikemodel "food-delivery/module/restaurantlike/model"
	"food-delivery/pubsub"
	"log"
)

type userDislikeRestaurantStore interface {
	//FindDataWithConditions(
	//	context context.Context,
	//	condition map[string]interface{},
	//) (*restaurantlikemodel.Like, error)
	Delete(context context.Context, userId int, restaurantId int) error
}

//type DecreaseLikeCountRstStore interface {
//	DecreasedLikeCount(context context.Context, id int) error
//}

type userDislikeRestaurantBiz struct {
	store userDislikeRestaurantStore
	//decStore DecreaseLikeCountRstStore
	ps pubsub.Pubsub
}

func NewUserDislikeRestaurantBiz(store userDislikeRestaurantStore, ps pubsub.Pubsub) *userDislikeRestaurantBiz {
	return &userDislikeRestaurantBiz{
		store: store,
		//decStore: decStore,
		ps: ps,
	}
}

func (biz *userDislikeRestaurantBiz) DislikeRestaurant(
	ctx context.Context,
	userId,
	restaurantId int,
) error {
	err := biz.store.Delete(ctx, userId, restaurantId)

	if err != nil {
		return restaurantlikemodel.ErrCannotLikeRestaurant(err)
	}

	//send message
	if err := biz.ps.Publish(ctx, common.TopicUserDislikeRestaurant, pubsub.NewMessage(&restaurantlikemodel.Like{RestaurantId: restaurantId, UserId: userId})); err != nil {
		log.Println(err)
	}

	//j := asyncjob.NewJob(func(ctx context.Context) error {
	//	return biz.decStore.DecreasedLikeCount(ctx, restaurantId)
	//})
	//
	//asyncjob.NewGroup(true, j).Run(ctx)

	return nil
}
