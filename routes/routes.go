package routes

import (
	"sensibull/stocks-api/business/router"
	"sensibull/stocks-api/middleware"
	"sensibull/stocks-api/middleware/corel"
	"sensibull/stocks-api/utils/logging"

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
