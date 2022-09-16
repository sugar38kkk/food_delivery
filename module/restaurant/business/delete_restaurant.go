package restaurantbusiness

import (
	"context"
	"errors"
	restaurantmodel "food-delivery/module/restaurant/model"
)

type DeleteRestaurantStore interface {
	FindDataWithCondition(
		context context.Context,
		condition map[string]interface{},
		moreKeys ...string,
	) (*restaurantmodel.Restaurant, error)
	Delete(context context.Context, id int) error
}

type deleteRestaurantBusiness struct {
	store DeleteRestaurantStore
}

func NewDeleteRestaurantBusiness(store DeleteRestaurantStore) *deleteRestaurantBusiness {
	return &deleteRestaurantBusiness{store: store}
}

func (biz *deleteRestaurantBusiness) DeleteRestaurant(context context.Context, id int) error {
	oldData, err := biz.store.FindDataWithCondition(context, map[string]interface{}{"id": id})

	if err != nil {
		return nil
	}

	if oldData.Status == 0 {
		return errors.New("Data has been deleted.")
	}

	if err := biz.store.Delete(context, id); err != nil {
		return err
	}

	return nil
}
