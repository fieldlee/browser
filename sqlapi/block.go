package sqlapi

import (
	"browser/handle"
	"browser/model"
	"browser/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)
// Get INfo
func GetInfo (c *gin.Context) {
	fabsdk := handle.InitSdk()
	defer fabsdk.Close()
	response,err := fabsdk.GetInfo()
	if err != nil {
		c.JSON(http.StatusInternalServerError,gin.H{
			"success":false,
			"err":err.Error(),
		})
		return
	}
	response.BCI.Height = response.BCI.Height - 1
	c.JSON(http.StatusOK,gin.H{
		"success":true,
		"info":response,
	})
}
// Get Block by height
func GetBlocksByHeight(c *gin.Context) {
	fabsdk := handle.InitSdk()
	defer fabsdk.Close()
	strStart := c.Param("start")
	strLimit := c.Param("limit")
	start , err := strconv.Atoi(strStart)
	if err != nil {
		c.JSON(http.StatusInternalServerError,gin.H{
			"success":false,
			"err":err.Error(),
		})
		return
	}
	limit , err := strconv.Atoi(strLimit)
	if err != nil {
		c.JSON(http.StatusInternalServerError,gin.H{
			"success":false,
			"err":err.Error(),
		})
		return
	}

	listBlocks := make([]model.Block,0)

	sqlClient,err := utils.InitSql()
	if err != nil {
		c.JSON(http.StatusInternalServerError,gin.H{
			"success":false,
			"err":err.Error(),
		})
		return
	}
	defer sqlClient.CloseSql()

	listBh,err := sqlClient.QueryBlocksByRange(start,limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError,gin.H{
			"success":false,
			"err":err.Error(),
		})
		return
	}

	for i := 0;i < len(listBh);i++{

		txs , err := sqlClient.QueryTxsByBlockHash(listBh[i].DataHash)
		if err != nil {
			fmt.Printf("Query Txs By block hash err :%s",err.Error())
			continue
		}


		tmpBck := model.Block{
			Number:listBh[i].Number,
			PreviousHash:listBh[i].PreviousHash,
			DataHash:listBh[i].DataHash,
			TxList:txs,
		}

		listBlocks = append(listBlocks,tmpBck)
	}

	c.JSON(http.StatusOK,gin.H{
		"success":true,
		"blocks":listBlocks,
	})
	return
}
// Get Block BY height
func GetBlockByHeight(c *gin.Context) {
	fabsdk := handle.InitSdk()
	defer fabsdk.Close()
	strHeight := c.Param("height")
	height , err := strconv.Atoi(strHeight)
	if err != nil {
		c.JSON(http.StatusInternalServerError,gin.H{
			"success":false,
			"err":err.Error(),
		})
		return
	}
	sqlClient,err := utils.InitSql()
	if err != nil {
		c.JSON(http.StatusInternalServerError,gin.H{
			"success":false,
			"err":err.Error(),
		})
		return
	}
	defer sqlClient.CloseSql()

	bh,err := sqlClient.QueryBlockByHeight(height)
	if err != nil {
		c.JSON(http.StatusInternalServerError,gin.H{
			"success":false,
			"err":err.Error(),
		})
		return
	}

	txs , err := sqlClient.QueryTxsByBlockHash(bh.DataHash)
	if err != nil {
		c.JSON(http.StatusInternalServerError,gin.H{
			"success":false,
			"err":err.Error(),})
		return
	}
	tmpBck := model.Block{
		Number:bh.Number,
		PreviousHash:bh.PreviousHash,
		DataHash:bh.DataHash,
		TxList:txs,
	}
	c.JSON(http.StatusOK,gin.H{
		"success":true,
		"block":tmpBck,
	})
	return
}
// Get BLock by hash
func GetBlockByHash(c *gin.Context) {
	fabsdk := handle.InitSdk()
	defer fabsdk.Close()
	hash := c.Param("hash")
	sqlClient,err := utils.InitSql()
	if err != nil {
		c.JSON(http.StatusInternalServerError,gin.H{
			"success":false,
			"err":err.Error(),
		})
		return
	}
	defer sqlClient.CloseSql()
	bh,err := sqlClient.QueryBlockByHash(hash)
	if err != nil {
		c.JSON(http.StatusInternalServerError,gin.H{
			"success":false,
			"err":err.Error(),
		})
		return
	}
	txs , err := sqlClient.QueryTxsByBlockHash(bh.DataHash)
	if err != nil {
		c.JSON(http.StatusInternalServerError,gin.H{
			"success":false,
			"err":err.Error(),})
		return
	}
	tmpBck := model.Block{
		Number:bh.Number,
		PreviousHash:bh.PreviousHash,
		DataHash:bh.DataHash,
		TxList:txs,
	}
	c.JSON(http.StatusOK,gin.H{
		"success":true,
		"block":tmpBck,
	})
	return
}
// Get Tx by Txid
func GetTxByID(c *gin.Context) {
	fabsdk := handle.InitSdk()
	defer fabsdk.Close()
	hash := c.Param("id")

	sqlClient,err := utils.InitSql()
	if err != nil {
		c.JSON(http.StatusInternalServerError,gin.H{
			"success":false,
			"err":err.Error(),
		})
		return
	}
	defer sqlClient.CloseSql()
	txinfo,err := sqlClient.QueryTxs(hash)
	if err != nil {
		c.JSON(http.StatusInternalServerError,gin.H{
			"success":false,
			"err":err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK,gin.H{
		"success":true,
		"tx":txinfo,
	})
	return
}
// Get BLock by hash
func GetBlockByTxHash(c *gin.Context) {
	fabsdk := handle.InitSdk()
	defer fabsdk.Close()
	hash := c.Param("hash")
	sqlClient,err := utils.InitSql()
	if err != nil {
		c.JSON(http.StatusInternalServerError,gin.H{
			"success":false,
			"err":err.Error(),
		})
		return
	}
	defer sqlClient.CloseSql()
	blockhash,err := sqlClient.QueryBlockHashByTxId(hash)
	if err != nil {
		c.JSON(http.StatusInternalServerError,gin.H{
			"success":false,
			"err":err.Error(),
		})
		return
	}

	bh,err := sqlClient.QueryBlockByHash(blockhash)

	if err != nil {
		c.JSON(http.StatusInternalServerError,gin.H{
			"success":false,
			"err":err.Error(),
		})
		return
	}
	txs , err := sqlClient.QueryTxsByBlockHash(bh.DataHash)
	if err != nil {
		c.JSON(http.StatusInternalServerError,gin.H{
			"success":false,
			"err":err.Error(),})
		return
	}
	tmpBck := model.Block{
		Number:bh.Number,
		PreviousHash:bh.PreviousHash,
		DataHash:bh.DataHash,
		TxList:txs,
	}
	c.JSON(http.StatusOK,gin.H{
		"success":true,
		"block":tmpBck,
	})
	return
}

// Get Tx by account
func GetTxsByAccount(c *gin.Context){
	fabsdk := handle.InitSdk()
	defer fabsdk.Close()
	account := c.Param("account")

	sqlClient,err := utils.InitSql()
	if err != nil {
		c.JSON(http.StatusInternalServerError,gin.H{
			"success":false,
			"err":err.Error(),
		})
		return
	}
	defer sqlClient.CloseSql()
	txinfo,err := sqlClient.QueryTxsByAccount(account)
	if err != nil {
		c.JSON(http.StatusInternalServerError,gin.H{
			"success":false,
			"err":err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK,gin.H{
		"success":true,
		"txs":txinfo,
	})
	return
}

// Get Txs by token
func GetTxsByToken(c *gin.Context){
	fabsdk := handle.InitSdk()
	defer fabsdk.Close()
	token := c.Param("token")

	sqlClient,err := utils.InitSql()
	if err != nil {
		c.JSON(http.StatusInternalServerError,gin.H{
			"success":false,
			"err":err.Error(),
		})
		return
	}
	defer sqlClient.CloseSql()
	txinfo,err := sqlClient.QueryTxsByToken(token)
	if err != nil {
		c.JSON(http.StatusInternalServerError,gin.H{
			"success":false,
			"err":err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK,gin.H{
		"success":true,
		"txs":txinfo,
	})
	return
}

// Get Txs heigth
func GetTxHeight(c *gin.Context){
	sqlClient,err := utils.InitSql()
	if err != nil {
		c.JSON(http.StatusInternalServerError,gin.H{
			"success":false,
			"err":err.Error(),
		})
		return
	}
	defer sqlClient.CloseSql()

	count,err := sqlClient.QueryTxsNum()
	if err != nil {
		c.JSON(http.StatusInternalServerError,gin.H{
			"success":false,
			"err":err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK,gin.H{
		"success":true,
		"txcount":count,
	})
	return
}
// Get Txs heigth
func GetTxsByHeigth(c *gin.Context){
	sqlClient,err := utils.InitSql()
	if err != nil {
		c.JSON(http.StatusInternalServerError,gin.H{
			"success":false,
			"err":err.Error(),
		})
		return
	}
	defer sqlClient.CloseSql()

	strStart := c.Param("start")
	strLimit := c.Param("limit")
	start , err := strconv.Atoi(strStart)
	if err != nil {
		c.JSON(http.StatusInternalServerError,gin.H{
			"success":false,
			"err":err.Error(),
		})
		return
	}
	limit , err := strconv.Atoi(strLimit)
	if err != nil {
		c.JSON(http.StatusInternalServerError,gin.H{
			"success":false,
			"err":err.Error(),
		})
		return
	}

	txs,err :=sqlClient.QueryTxsByRange(start,limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError,gin.H{
			"success":false,
			"err":err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK,gin.H{
		"success":true,
		"txs":txs,
	})
	return
}