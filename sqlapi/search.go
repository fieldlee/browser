package sqlapi

import (
	"browser/utils"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func Search(c *gin.Context){
	type SearchParam struct {
		Search string   `json:"search"`
	}
	param := new(SearchParam)
	err := c.BindJSON(param)
	if err != nil {
		c.JSON(http.StatusInternalServerError,gin.H{
			"success":false,
			"err":err.Error(),
		})
		return
	}
	key := param.Search

	sqlClient,err := utils.InitSql()
	if err != nil {
		c.JSON(http.StatusInternalServerError,gin.H{
			"success":false,
			"err":"init error",
		})
		return
	}
	defer sqlClient.CloseSql()

	// block
	blockHeader,err := sqlClient.QueryBlockByHash(key)
	if err != nil && err != sql.ErrNoRows {
		c.JSON(http.StatusInternalServerError,gin.H{
			"success": false,
			"err":     fmt.Sprintf("get block by hash err :%s", err.Error()),
		})
		return
	}
	if blockHeader.DataHash != "" {
		//return
		c.JSON(http.StatusOK,gin.H{
			"success":true,
			"type":"block",
			"info":blockHeader,
		})
		return
	}


	intkey,err := strconv.Atoi(key)
	if err == nil {

		block,err := sqlClient.QueryBlockByHeight(intkey)
		if err != nil && err != sql.ErrNoRows {
			c.JSON(http.StatusInternalServerError,gin.H{
				"success":false,
				"err":fmt.Sprintf("query block by hash err :%s",err.Error()),
			})
			return
		}
		if block.DataHash != ""{
			//return
			c.JSON(http.StatusOK,gin.H{
				"success":true,
				"type":"block",
				"info":block,
			})
			return
		}

	}


	// tx
	txs,err := sqlClient.QueryTxs(key)
	if err != nil && err != sql.ErrNoRows {
		c.JSON(http.StatusInternalServerError,gin.H{
			"success":false,
			"err":fmt.Sprintf("query transactions by hash err :%s",err.Error()),
		})
		return
	}
	if txs.TransactionId != "" {
		//return
		c.JSON(http.StatusOK,gin.H{
			"success":true,
			"type":"transaction",
			"info":txs,
		})
		return
	}


	// token
	tokens,err := sqlClient.QueryTokensById(key)
	if err != nil && err != sql.ErrNoRows {
		c.JSON(http.StatusInternalServerError,gin.H{
			"success":false,
			"err":fmt.Sprintf("query token by id err :%s",err.Error()),
		})
		return
	}
	if len(tokens)>0{
		//return
		c.JSON(http.StatusOK,gin.H{
			"success":true,
			"type":"token",
			"info":tokens,
		})
		return
	}

	// account
	couchClient,err := utils.InitCouchClient()
	if err != nil  {
		c.JSON(http.StatusInternalServerError,gin.H{
			"success":false,
			"err":err.Error(),
		})
		return
	}
	accountList,err := couchClient.GetAccount(key)
	if err != nil {
		c.JSON(http.StatusInternalServerError,gin.H{
			"success":false,
			"err":err.Error(),
		})
		return
	}

	if accountList.Name != ""{
		//return
		c.JSON(http.StatusOK,gin.H{
			"success":true,
			"type":"account",
			"info":accountList,
		})
		return
	}
	c.JSON(http.StatusOK,gin.H{
		"success":true,
		"type":"false",
	})
	return
}