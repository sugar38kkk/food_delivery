package userbiz

import (
	"context"
	"food-delivery/common"
	"food-delivery/component/appctx"
	"food-delivery/component/tokenprovider"
	usermodel "food-delivery/module/user/model"
)

type LoginStorage interface {
	FindDataWithCondition(context context.Context, condition map[string]interface{}, moreKeys ...string) (*usermodel.User, error)
}

type LoginBiz struct {
	appCtx        appctx.AppContext
	store         LoginStorage
	tokenProvider tokenprovider.Provider
	hasher        Hasher
	expiry        int
}

func NewLoginBiz(store LoginStorage, hasher Hasher, tokenProvider tokenprovider.Provider, expiry int) *LoginBiz {
	return &LoginBiz{
		store:         store,
		hasher:        hasher,
		tokenProvider: tokenProvider,
		expiry:        expiry,
	}
}

func (biz *LoginBiz) Login(ctx context.Context, data *usermodel.UserLogin) (*tokenprovider.Token, error) {
	user, err := biz.store.FindDataWithCondition(ctx, map[string]interface{}{"email": data.Email})

	if err != nil {
		return nil, usermodel.ErrEmailOrPasswordInvalid
	}

	passHashed := biz.hasher.Hash(data.Password + user.Salt)

	if user.Password != passHashed {
		return nil, usermodel.ErrEmailOrPasswordInvalid
	}

	payload := tokenprovider.TokenPayload{
		UserId: user.Id,
		Role:   user.Role,
	}

	accessToken, err := biz.tokenProvider.Generate(payload, biz.expiry)
	if err != nil {
		return nil, common.ErrInternal(err)
	}

	return accessToken, nil
}
