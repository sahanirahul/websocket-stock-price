package router

import (
	corehttp "net/http"
	"sensibull/stocks-api/business/interfaces/iusecase"
	"sensibull/stocks-api/business/repository/db"
	"sensibull/stocks-api/business/repository/http"
	"sensibull/stocks-api/business/repository/websocket"
	"sensibull/stocks-api/business/usecase"
	"sync"

	"github.com/gin-gonic/gin"
)

func provideInstrumentRouter() *instrumentRouter {
	instrumentService := usecase.NewInstrumentService(http.NewInstrumentHttpRepo(), websocket.NewWebsocketRepo(), db.NewInstrumentRepo())
	return newInstrumentRouter(instrumentService)

}

func InstrumentRoutes(apigroup *gin.RouterGroup) {
	r := provideInstrumentRouter()
	apigroup.GET("underlying-prices", r.getUnderlyingPrices)
}

type instrumentRouter struct {
	instrumentService iusecase.IStocksInstrumentsService
}

var routeOnce sync.Once
var insrouter *instrumentRouter

func newInstrumentRouter(is iusecase.IStocksInstrumentsService) *instrumentRouter {
	routeOnce.Do(func() {
		insrouter = &instrumentRouter{
			instrumentService: is,
		}
	})
	return insrouter
}

func (ir *instrumentRouter) getUnderlyingPrices(c *gin.Context) {
	c.JSON(corehttp.StatusOK, gin.H{
		"success": true,
		"payload": nil,
	})
}
