package api

import (
	"browser/handle"
	"browser/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SyncToken(c *gin.Context) {
	fabsdk := handle.InitSdk()
	defer fabsdk.Close()
	tokens,err := fabsdk.GetTokens()
	if err != nil{
		c.JSON(http.StatusInternalServerError,gin.H{
			"err":err.Error(),
		})
		return
	}
	sqlClient,err := utils.InitSql()
	if err != nil {
		c.JSON(http.StatusInternalServerError,gin.H{
			"err":"init error",
		})
		return
	}
	defer sqlClient.CloseSql()
	for i := 0;i<len(tokens) ;i++  {
		err = sqlClient.InsertToken(tokens[i])
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
	}
	c.JSON(http.StatusOK,gin.H{
		"tokens":tokens,
	})
	return
}