package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func InitRoutes(router *gin.Engine) {
	api := router.Group("/api")
	publicGroup := api.Group("/public")
	publicGroup.GET("v1/test", func(c *gin.Context) { c.JSON(http.StatusOK, "OK") })

}
