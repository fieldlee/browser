package sqlapi

import (
	"browser/handle"
	"browser/model"
	"browser/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Token(c *gin.Context) {
	sqlClient,err := utils.InitSql()
	if err != nil {
		c.JSON(http.StatusInternalServerError,gin.H{
			"success":false,
			"err":"init error",
		})
		return
	}
	defer sqlClient.CloseSql()

	tokens , err := sqlClient.QueryTokens()

	if err != nil {
		c.JSON(http.StatusInternalServerError,gin.H{
			"success":false,
			"err":err.Error(),
		})
		return
	}

	tList := []model.TokenInfo{}

	for _, t := range tokens {
		ti := model.TokenInfo{}
		ti.Name = t.Name
		ti.Status = t.Status
		ti.Amount = t.Amount
		ti.Desc = t.Desc
		ti.Issuer = t.Issuer
		txs,err := sqlClient.QueryTxsByToken(ti.Name)
		if err != nil {
			ti.TxsNumber = 0
		}else {
			ti.TxsNumber = len(txs)
		}
		tList = append(tList,ti)
	}

	c.JSON(http.StatusOK,gin.H{
		"success":true,
		"info":tList,
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
			"success":false,
			"err":err.Error(),
		})
		return
	}
	hisList,err := fabsdk.GetTokenHistory(postToken.Token)
	if err != nil {
		c.JSON(http.StatusInternalServerError,gin.H{
			"success":false,
			"err":err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK,gin.H{
		"success":true,
		"info":hisList,
	})
	return
}