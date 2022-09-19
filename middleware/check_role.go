package middleware

import (
	"errors"
	"food-delivery/common"
	"food-delivery/component/appctx"
	"github.com/gin-gonic/gin"
)

func CheckRole(appCtx appctx.AppContext, allowRoles ...string) func(c *gin.Context) {
	return func(c *gin.Context) {
		u := c.MustGet(common.CurrentUser).(common.Requester)

		isValid := false

		for _, item := range allowRoles {
			if u.GetRole() == item {
				isValid = true
				break
			}
		}

		if !isValid {
			panic(common.ErrNoPermission(errors.New("invalid role user")))
		}

		c.Next()

	}
}
