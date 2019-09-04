package api

import (
	"browser/handle"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Invoke(c *gin.Context){
	type AccountParam struct {
		Fcn string    `json:"fcn"`
		Args []string `json:"args"`
	}
	param := new(AccountParam)
	err := c.BindJSON(param)
	if err != nil {
		c.JSON(http.StatusInternalServerError,gin.H{
			"success":false,
			"err":err.Error(),
		})
		return
	}

	fasdk := handle.InitSdk()
	defer fasdk.Close()
	message,txid,err := fasdk.Invoke(param.Fcn,param.Args)
	if err != nil {
		c.JSON(http.StatusInternalServerError,gin.H{
			"success":false,
			"err":err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK,gin.H{
		"success":true,
		"transactionId":txid,
		"info":message,
	})
	return
}

func Query(c *gin.Context){
	type AccountParam struct {
		Fcn string    `json:"fcn"`
		Args []string `json:"args"`
	}
	param := new(AccountParam)
	err := c.BindJSON(param)
	if err != nil {
		c.JSON(http.StatusInternalServerError,gin.H{
			"success":false,
			"err":err.Error(),
		})
		return
	}

	fasdk := handle.InitSdk()
	defer fasdk.Close()
	message,err := fasdk.Query(param.Fcn,param.Args)
	if err != nil {
		c.JSON(http.StatusInternalServerError,gin.H{
			"success":false,
			"err":err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK,gin.H{
		"success":true,
		"info":message,
	})
	return
}