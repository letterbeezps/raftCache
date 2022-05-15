package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/letterbeezps/raftCache/global"
	"go.uber.org/zap"
)

type JoinRequest struct {
	Id   string `json:"id"`
	Addr string `json:"addr"`
}

func TestGet(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"success": true,
		"msg":     "testGet",
	})
}

func Get(ctx *gin.Context) {
	key := ctx.Param("key")

	value, err := global.RaftServer.Get(key)

	if err != nil {
		zap.S().Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"err": err.Error(),
		})
	}

	res := make(map[string]interface{})

	json.Unmarshal(value, &res)

	ctx.JSON(http.StatusOK, gin.H{
		"key":   key,
		"value": res,
	})
}

func Delete(ctx *gin.Context) {
	key := ctx.Param("key")

	err := global.RaftServer.Delete(key)

	if err != nil {
		zap.S().Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"err": err.Error(),
		})
	}

	ctx.JSON(http.StatusOK, "")
}

func Set(ctx *gin.Context) {
	key := ctx.Param("key")

	value, err := ioutil.ReadAll(ctx.Request.Body)

	defer ctx.Request.Body.Close()

	err = global.RaftServer.Set(key, value)

	if err != nil {
		zap.S().Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"err": err.Error(),
		})
	}

	ctx.JSON(http.StatusOK, "")
}

func Join(ctx *gin.Context) {
	zap.S().Infof("get join request from %s", ctx.Request.RequestURI)

	joinReq := JoinRequest{}

	ctx.BindJSON(&joinReq)

	if joinReq.Addr == "" {
		zap.S().Error("joibAddr can not be null")
	}

	if joinReq.Id == "" {
		zap.S().Error("joibId can not be null")
	}

	err := global.RaftServer.Join(joinReq.Id, joinReq.Addr)

	if err != nil {
		zap.S().Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"err": err.Error(),
		})
	}

	ctx.JSON(http.StatusOK, "")
}
