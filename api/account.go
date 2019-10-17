package api

import (
	"browser/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetAccounts(c *gin.Context) {
	couchClient,err := utils.InitCouchClient()
	if err != nil {
		c.JSON(http.StatusInternalServerError,gin.H{
			"success":false,
			"err":err.Error(),
		})
		return
	}
	accountList,err := couchClient.GetAccounts()
	if err != nil {
		c.JSON(http.StatusInternalServerError,gin.H{
			"success":false,
			"err":err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK,gin.H{
		"success":true,
		"info":accountList,
	})
	return
}

func GetAccount(c *gin.Context) {
	couchClient,err := utils.InitCouchClient()
	if err != nil {
		c.JSON(http.StatusInternalServerError,gin.H{
			"success":false,
			"err":err.Error(),
		})
		return
	}

	type Account struct {
		Name string `json:"name"`
	}
	account := new(Account)
	err = c.BindJSON(account)
	if err != nil {
		c.JSON(http.StatusInternalServerError,gin.H{
			"success":false,
			"err":err.Error(),
		})
		return
	}
	accountList,err := couchClient.GetAccount(account.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError,gin.H{
			"success":false,
			"err":err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK,gin.H{
		"success":true,
		"info":accountList,
	})
	return
}
