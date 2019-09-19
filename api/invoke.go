package api

import (
	"browser/handle"
	"browser/model"
	"encoding/json"
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

func QueryHold(c *gin.Context){
	fasdk := handle.InitSdk()
	defer fasdk.Close()
	account := c.Param("account")
	message,err := fasdk.Query("holdtoken",[]string{account})
	if err != nil {
		c.JSON(http.StatusInternalServerError,gin.H{
			"success":false,
			"err":err.Error(),
		})
		return
	}

	var actions = []model.LedgerAction{}
	err = json.Unmarshal([]byte(message.Message),&actions)
	if err != nil {
		c.JSON(http.StatusInternalServerError,gin.H{
			"success":false,
			"err":err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK,gin.H{
		"success":true,
		"info":actions,
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