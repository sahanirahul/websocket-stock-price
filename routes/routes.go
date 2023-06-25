package routes

import (
	"sensibull/stocks-api/business/router"

	"github.com/gin-gonic/gin"
)

func InitRoutes(baserouter *gin.Engine) {
	api := baserouter.Group("")
	router.InstrumentRoutes(api)
}
