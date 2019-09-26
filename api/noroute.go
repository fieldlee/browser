package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func NoRouterHandle(c *gin.Context){
	c.JSON(http.StatusNotFound,gin.H{
		"success":false,
		"msg":"the router not found",
	})
}
