package router

import (
	"context"
	"fmt"
	"log"
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
	err := instrumentService.UpdateEquityStockDetails(context.Background())
	if err != nil {
		fmt.Println("unable to run initial stock update job")
		log.Fatal(err)
	}
	go func() {
		err := instrumentService.UpdateDerivativeStockDetails(context.Background())
		if err != nil {
			fmt.Println("unable to run initial derivative update job")
			log.Fatal(err)
		}
	}()

	return newInstrumentRouter(instrumentService)

}

func InstrumentRoutes(apigroup *gin.RouterGroup) {
	r := provideInstrumentRouter()
	apigroup.GET("underlying-prices", r.getUnderlyingPrices)
	apigroup.GET("derivative-prices/:symbol", r.getDerivativePrices)
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
	payload, err := ir.instrumentService.FetchEquityStockDetails(c)
	if err != nil {
		c.JSON(corehttp.StatusOK, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	c.JSON(corehttp.StatusOK, gin.H{
		"success": true,
		"payload": payload,
	})
}

func (ir *instrumentRouter) getDerivativePrices(c *gin.Context) {
	var req struct {
		Symbol string `uri:"symbol" binding:"required,gt=0"`
	}
	err := c.ShouldBindUri(&req)
	if err != nil {
		c.JSON(corehttp.StatusBadRequest, gin.H{
			"success": false,
			"error":   "invalid_request",
		})
		return
	}
	payload, err := ir.instrumentService.FetchDerivativeStockDetails(c, req.Symbol)
	if err != nil {
		c.JSON(corehttp.StatusUnprocessableEntity, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	c.JSON(corehttp.StatusOK, gin.H{
		"success": true,
		"payload": payload,
	})
}
