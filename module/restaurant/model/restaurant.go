package restaurantmodel

import (
	"errors"
	"food-delivery/common"
	"strings"
)

const TableName string = "restaurants"
const EntityName string = "Restaurants"

type Restaurant struct {
	common.SQLModel `json:",inline"`
	Name            string         `json:"name" gorm:"column:name"`
	Addr            string         `json:"addr" gorm:"column:addr"`
	Logo            *common.Image  `json:"logo" gorm:"logo"`
	Cover           *common.Images `json:"cover" gorm:"cover"`
}

func (Restaurant) TableName() string { return TableName }

func (r *Restaurant) Mask(isAdminOrOwner bool) {
	r.GenUID(common.DbTypeRestaurant)
}

type RestaurantCreate struct {
	common.SQLModel `json:",inline"`
	Name            string         `json:"name" gorm:"column:name"`
	Addr            string         `json:"addr" gorm:"column:addr"`
	Logo            *common.Image  `json:"logo" gorm:"logo"`
	Cover           *common.Images `json:"cover" gorm:"cover"`
}

func (data *RestaurantCreate) Validate() error {
	data.Name = strings.TrimSpace(data.Name)

	if data.Name == "" {
		return ErrNameIsEmpty
	}

	return nil
}

func (data *RestaurantCreate) Mask(isAdminOrOwner bool) {
	data.GenUID(common.DbTypeRestaurant)
}

func (RestaurantCreate) TableName() string { return TableName }

type RestaurantUpdate struct {
	Name  *string        `json:"name" gorm:"column:name"`
	Addr  *string        `json:"addr" gorm:"column:addr"`
	Logo  *common.Image  `json:"logo" gorm:"logo"`
	Cover *common.Images `json:"cover" gorm:"cover"`
}

func (RestaurantUpdate) TableName() string { return TableName }

var (
	ErrNameIsEmpty = errors.New("Name cannot be empty.")
)
