package v1

import (
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

	c.JSON(200, gin.H{
		"status": "ok",
	})
}
