package route

import (
	"github.com/gin-gonic/gin"
	"github.com/letterbeezps/raftCache/internal/api"
	"go.uber.org/zap"
)

func InitCacheRouter(Router *gin.RouterGroup) {
	CacheRouter := Router.Group("/cache")

	zap.S().Infof("register user router")

	{
		// CacheRouter.GET("/test", api.TestGet)
		CacheRouter.GET(":key", api.Get)
		CacheRouter.DELETE(":key", api.Delete)
		CacheRouter.POST(":key", api.Set)
		CacheRouter.POST("/join", api.Join)
	}
}
