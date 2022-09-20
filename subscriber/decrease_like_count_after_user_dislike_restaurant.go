package subscriber

import (
	"context"
	"food-delivery/component/appctx"
	restaurantstorage "food-delivery/module/restaurant/storage"
	"food-delivery/pubsub"
)

func DecreaseLikeCountAfterUserDislikeRestaurant(appCtx appctx.AppContext) consumerJob {
	return consumerJob{
		Title: "Decrease like count after user dislikes restaurant",

		Hld: func(ctx context.Context, message *pubsub.Message) error {

			store := restaurantstorage.NewSQLStore(appCtx.GetMainDBConnection())

			likeData := message.Data().(HasRestaurantId)

			return store.DecreasedLikeCount(ctx, likeData.GetRestaurantId())
		},
	}
}
