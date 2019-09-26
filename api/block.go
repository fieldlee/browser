package api

import (
	"browser/handle"
	"browser/model"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)
// Get blocks
func GetBlocks (c *gin.Context) {
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
	curHeight := response.BCI.Height
	listBlocks := make([]model.Block,0)
	for i:=0;i<10 ;i++  {
		t := curHeight - uint64(i)
		if t >= 0{
			if blockinfo ,err := fabsdk.GetBlocks(t);err != nil {
				fmt.Errorf(err.Error())
				continue
			}else{
				listBlocks = append(listBlocks,blockinfo)
			}
		}
	}
	c.JSON(http.StatusOK,gin.H{
		"success":true,
		"blocks":listBlocks,
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
	curHeight := uint64(start)
	listBlocks := make([]model.Block,0)
	for i:=0;i<limit ;i++  {
		t := curHeight - uint64(i)
		if t >= 0{
			if blockinfo ,err := fabsdk.GetBlocks(t);err != nil {
				continue
			}else{
				listBlocks = append(listBlocks,blockinfo)
			}
		}
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
	blockinfo,err := fabsdk.GetBlocks(uint64(height))
	if err != nil {
		c.JSON(http.StatusInternalServerError,gin.H{
			"success":false,
			"err":err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK,gin.H{
		"success":true,
		"block":blockinfo,
	})
	return
}

// Get BLock by hash
func GetBlockByHash(c *gin.Context) {
	fabsdk := handle.InitSdk()
	defer fabsdk.Close()
	hash := c.Param("hash")
	blockinfo,err := fabsdk.GetBlocksByHash(hash)
	if err != nil {
		c.JSON(http.StatusInternalServerError,gin.H{
			"success":false,
			"err":err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK,gin.H{
		"success":true,
		"block":blockinfo,
	})
	return
}

// Get Tx by Txid
func GetTxByID(c *gin.Context) {
	fabsdk := handle.InitSdk()
	defer fabsdk.Close()
	hash := c.Param("id")
	txinfo,err := fabsdk.GetTransactionByTxId(hash)
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
	blockinfo,err := fabsdk.GetBlocksByTxId(hash)
	if err != nil {
		c.JSON(http.StatusInternalServerError,gin.H{
			"success":false,
			"err":err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK,gin.H{
		"success":true,
		"block":blockinfo,
	})
	return
}



