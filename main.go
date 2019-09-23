package main

import (
	"browser/api"
	"browser/sqlapi"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

func Logger() gin.HandlerFunc {
	logClient := logrus.New()
	logClient.SetLevel(logrus.DebugLevel)
	return func (c *gin.Context) {
		// 开始时间
		start := time.Now()
		// 处理请求
		c.Next()
		// 结束时间
		end := time.Now()
		//执行时间
		latency := end.Sub(start)
		path := c.Request.URL.Path
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		logClient.Infof("| %3d | %13v | %15s | %s  %s |",
			statusCode,
			latency,
			clientIP,
			method, path,
		)
	}
}

func setupRouter() *gin.Engine {
	gin.SetMode(gin.DebugMode)
	r := gin.Default()
	r.NoRoute(api.NoRouterHandle)
	r.Use(Logger())
	r.GET("/",func(c *gin.Context){
		c.Redirect(http.StatusMovedPermanently,"/s/info")
	})
	hlc := r.Group("/c")
	{
		block := hlc.Group("/block")
		{
			// 默认获得块列表
			block.GET("/list/:start/:limit", api.GetBlocksByHeight)
			// 默认获得制定高度的块
			block.GET("/height/:height", api.GetBlockByHeight)
			// 默认获得指定hash的块
			block.GET("/hash/:hash", api.GetBlockByHash)
			// 根据tx hash获得块
			block.GET("/tx/:hash", api.GetBlockByTxHash)
		}

		tx := hlc.Group("/tx")
		{
			// 默认获得指定hash的交易
			tx.GET("/id/:id",api.GetTxByID)
		}

		// 获得token
		hlc.GET("/token", api.SyncToken)
		// 获得account
		hlc.GET("/account", api.GetAccounts)

		// get account
		hlc.GET("/query/:account", api.QueryHold)
	}

	sql := r.Group("/s")
	{
		sqlblock:=sql.Group("/block")
		{
			// 默认获得块列表
			sqlblock.GET("/list/:start/:limit", sqlapi.GetBlocksByHeight)
			// 默认获得制定高度的块
			sqlblock.GET("/height/:height", sqlapi.GetBlockByHeight)
			// 默认获得指定hash的块
			sqlblock.GET("/hash/:hash", sqlapi.GetBlockByHash)
			// 根据tx hash获得块
			sqlblock.GET("/tx/:hash", sqlapi.GetBlockByTxHash)
		}
		sqltx := sql.Group("/tx")
		{
			sqltx.GET("/id/:id",sqlapi.GetTxByID)
		}
		sql.GET("/info",sqlapi.GetInfo)
		sql.GET("/token",sqlapi.Token)
		sql.GET("/token/:token",sqlapi.GetTxsByToken)
		sql.POST("/token",sqlapi.TokenHistory)
		sql.POST("/account",api.GetAccount)
		sql.GET("/txs/:account",sqlapi.GetTxsByAccount)
	}

	r.POST("/invoke",api.Invoke)
	r.POST("/query",api.Query)
	r.POST("/search",sqlapi.Search)
	// 同步block transaction
	r.GET("/sync", api.SyncBlockTx)
	r.GET("/synctoken", api.SyncToken)

	return r
}

func main() {
	//启动监听区块
	go api.ListenBlock()
	//启动service
	r := setupRouter()
	r.Run(":8080")
}
