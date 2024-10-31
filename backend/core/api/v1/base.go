package v1

import (
	"doo-store/backend/constant"
	"errors"

	"github.com/gin-gonic/gin"
)

func checkAuth(c *gin.Context, admin bool) error {
	token, tokenExist := c.Get("token")
	if !tokenExist {
		return errors.New(constant.ErrNoPermission)
	}
	t := token.(string)
	info, err := dootaskService.GetUserInfo(t)
	if err != nil {
		return err
	}
	if admin && !info.IsAdmin() {
		return errors.New(constant.ErrNoPermission)
	}
	return nil
}
