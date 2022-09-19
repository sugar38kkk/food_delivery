package restaurantbusiness

import (
	"context"
	"food-delivery/common"
	restaurantmodel "food-delivery/module/restaurant/model"
)

type ListRestaurantRepo interface {
	ListRestaurant(
		context context.Context,
		filter *restaurantmodel.Filter,
		paging *common.Paging,
	) ([]restaurantmodel.Restaurant, error)
}

type listRestaurantBusiness struct {
	repo ListRestaurantRepo
}

func NewListRestaurantBusiness(repo ListRestaurantRepo) *listRestaurantBusiness {
	return &listRestaurantBusiness{repo}
}

func (biz *listRestaurantBusiness) ListRestaurant(
	context context.Context,
	filter *restaurantmodel.Filter,
	paging *common.Paging,
) ([]restaurantmodel.Restaurant, error) {

	result, err := biz.repo.ListRestaurant(context, filter, paging)
	if err != nil {
		return nil, common.ErrCannotGetEntity(restaurantmodel.EntityName, err)
	}

	return result, nil
}
