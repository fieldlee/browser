package api

import (
	"browser/handle"
	"browser/model"
	"browser/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SyncBlockTx(c *gin.Context) {
	fabsdk := handle.InitSdk()
	defer fabsdk.Close()
	response,err := fabsdk.GetInfo()
	if err != nil {
		c.JSON(http.StatusInternalServerError,gin.H{
			"err":err.Error(),
		})
		return
	}
	curHeight := response.BCI.Height
	var curBlock model.BlockHeader

	sqlClient,err := utils.InitSql()

	if err != nil {
		c.JSON(http.StatusInternalServerError,gin.H{
			"err":"init error",
		})
		return
	}
	defer sqlClient.CloseSql()

	sqlHeight,_ := sqlClient.QueryBlockHeight()
	//if err != nil {
	//	c.JSON(http.StatusInternalServerError,gin.H{
	//		"err":err.Error(),
	//	})
	//	return
	//}
	start := 0
	if sqlHeight > 0 {
		start = sqlHeight + 1
	}
	for i := start; uint64(i) <= curHeight; i++  {
		if blockinfo ,err := fabsdk.GetBlocks(uint64(i));err != nil {
			fmt.Errorf(err.Error())
			continue
		}else{
			curBlock = model.BlockHeader{}
			curBlock.Number = blockinfo.Number
			curBlock.PreviousHash = blockinfo.PreviousHash
			curBlock.DataHash = blockinfo.DataHash
			// update createtime
			if len(blockinfo.TxList)>0{
				curBlock.CreateTime = blockinfo.TxList[0].CreateTime
			}

			err = sqlClient.InsertBlock(curBlock)
			if err != nil {
				fmt.Errorf(err.Error())
				continue
			}
			///update tx
			if len(blockinfo.TxList)>0{
				for _ , tx := range  blockinfo.TxList{
					if tx.TransactionId != ""{
						err = sqlClient.InsertTx(curBlock.DataHash,tx)
						if err != nil {
							fmt.Println(err.Error())
							continue
						}
					}
				}
			}
			//time.Sleep(time.Duration(500)*time.Microsecond)
			///update
			//updateBlock,err := sqlClient.QueryBlockByHeight(i-1)
			//if err != nil {
			//	fmt.Println(err.Error())
			//	continue
			//}
			//err = sqlClient.UpdateTxHash(curBlock.PreviousHash,updateBlock.DataHash)
			//if err != nil {
			//	fmt.Println(err.Error())
			//	continue
			//}
			//time.Sleep(1*time.Second)
			//err = sqlClient.UpdateBlockHash(i-1,curBlock.PreviousHash)
			//if err != nil {
			//	fmt.Println(err.Error())
			//	continue
			//}
		}
	}
	c.JSON(http.StatusOK,gin.H{
		"info":true,
	})
	return
}
