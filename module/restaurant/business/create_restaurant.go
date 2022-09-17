package restaurantbusiness

import (
	"context"
	"food-delivery/common"
	restaurantmodel "food-delivery/module/restaurant/model"
)

type CreateRestaurantStore interface {
	Create(context context.Context, data *restaurantmodel.RestaurantCreate) error
}

type createRestaurantBusiness struct {
	store CreateRestaurantStore
}

func NewCreateRestaurantBusiness(store CreateRestaurantStore) *createRestaurantBusiness {
	return &createRestaurantBusiness{store: store}
}

func (biz *createRestaurantBusiness) CreateRestaurant(context context.Context, data *restaurantmodel.RestaurantCreate) error {
	if err := data.Validate(); err != nil {
		return common.ErrInvalidRequest(err)
	}

	if err := biz.store.Create(context, data); err != nil {
		return common.ErrCannotCreateEntity(restaurantmodel.EntityName, err)
	}

	return nil
}
