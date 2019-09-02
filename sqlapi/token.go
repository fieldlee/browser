package sqlapi

import (
	"browser/handle"
	"browser/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Token(c *gin.Context) {
	sqlClient,err := utils.InitSql()
	if err != nil {
		c.JSON(http.StatusInternalServerError,gin.H{
			"err":"init error",
		})
		return
	}
	defer sqlClient.CloseSql()

	tokens , err := sqlClient.QueryTokens()

	if err != nil {
		c.JSON(http.StatusInternalServerError,gin.H{
			"err":err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK,gin.H{
		"tokens":tokens,
	})
	return
}

func TokenHistory(c *gin.Context){
	fabsdk := handle.InitSdk()
	defer fabsdk.Close()

	type PostToken struct{
		Token string `json:"token"`
	}

	var postToken PostToken

	err := c.BindJSON(&postToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError,gin.H{
			"err":err.Error(),
		})
		return
	}
	historylist,err :=fabsdk.GetTokenHistory(postToken.Token)
	if err != nil {
		c.JSON(http.StatusInternalServerError,gin.H{
			"err":err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK,gin.H{
		"history":historylist,
	})
	return
}