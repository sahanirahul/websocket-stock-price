package routes

import (
	"priceupdater/stocks-api/business/router"
	"priceupdater/stocks-api/middleware"
	"priceupdater/stocks-api/middleware/corel"
	"priceupdater/stocks-api/utils/logging"

	"github.com/gin-gonic/gin"
)

func InitRoutes(baserouter *gin.Engine) {
	api := baserouter.Group("")
	api.Use(corel.DefaultGinHandlers...)
	// adding recovery for api flow
	api.Use(middleware.Recovery(logging.Logger))
	api.Use(logging.Logger.Gin())
	router.InstrumentRoutes(api)
}
