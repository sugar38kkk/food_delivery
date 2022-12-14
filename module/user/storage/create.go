package userstorage

import (
	"context"
	"food-delivery/common"
	usermodel "food-delivery/module/user/model"
)

func (s *sqlStore) Create(context context.Context, data *usermodel.UserCreate) error {
	if err := s.db.Create(&data).Error; err != nil {
		return common.ErrDB(err)
	}
	return nil
}
