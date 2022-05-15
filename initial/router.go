package initial

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/letterbeezps/raftCache/internal/route"
)

func Routers() *gin.Engine {
	ApiRouter := gin.Default()

	ApiRouter.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"success": true,
		})
	})

	ApiGroupV1 := ApiRouter.Group("/v1")

	route.InitCacheRouter(ApiGroupV1)

	return ApiRouter
}
