package v1

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

// @Summary health
// @Schemes
// @Description
// @Tags public
// @Accept json
// @Produce json
// @Success 200 {string} string "ok"
// @Router /public/health [get]
func (b *BaseApi) HealthCheck(c *gin.Context) {
	token, ex1 := c.Get("token")
	userid, ex2 := c.Get("userid")
	fmt.Println(token, ex1, userid, ex2)
	c.JSON(200, gin.H{
		"status": "ok",
		"token":  token,
		"userid": userid,
	})
}
