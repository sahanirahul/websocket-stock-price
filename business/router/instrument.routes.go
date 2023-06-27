package router

import (
	"log"
	corehttp "net/http"
	"sensibull/stocks-api/business/interfaces/iusecase"
	"sensibull/stocks-api/business/repository/db"
	"sensibull/stocks-api/business/repository/http"
	"sensibull/stocks-api/business/repository/websocket"
	"sensibull/stocks-api/business/usecase"
	"sensibull/stocks-api/business/worker"
	"sensibull/stocks-api/middleware/corel"
	"sensibull/stocks-api/utils/logging"
	"sync"

	"github.com/gin-gonic/gin"
)

func provideInstrumentRouter() *instrumentRouter {
	dbrepo := db.NewInstrumentRepo()
	instrumentService := usecase.NewInstrumentService(http.NewInstrumentHttpRepo(), websocket.NewWebsocketRepo(dbrepo), dbrepo, worker.NewWorkerPool(50, 50))
	ctx := corel.CreateNewContext()
	logging.Logger.WriteLogs(ctx, "starting_initial_update_equity_details_job", logging.InfoLevel, logging.Fields{})
	_, err := instrumentService.UpdateEquityStockDetails(ctx)
	if err != nil {
		logging.Logger.WriteLogs(ctx, "unable_to_run_initial_equity_update_job", logging.ErrorLevel, logging.Fields{"error": err})
		log.Fatal(err)
	}
	go func() {
		logging.Logger.WriteLogs(ctx, "starting_initial_update_derivative_details_job", logging.InfoLevel, logging.Fields{})
		err := instrumentService.UpdateDerivativeStockDetails(ctx)
		if err != nil {
			logging.Logger.WriteLogs(ctx, "unable_to_run_initial_derivative_update_job", logging.ErrorLevel, logging.Fields{"error": err})
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
		logging.Logger.WriteLogs(c, "error_fetching_underlyings", logging.ErrorLevel, logging.Fields{"error": err})
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
		logging.Logger.WriteLogs(c, "error_fetching_underlying_derivatives", logging.ErrorLevel, logging.Fields{"error": err})
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
