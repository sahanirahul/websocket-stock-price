package router

import (
	"log"
	corehttp "net/http"
	"os"
	"priceupdater/stocks-api/business/interfaces/iusecase"
	"priceupdater/stocks-api/business/repository/db"
	"priceupdater/stocks-api/business/repository/http"
	"priceupdater/stocks-api/business/repository/websocket"
	"priceupdater/stocks-api/business/usecase"
	"priceupdater/stocks-api/business/worker"
	"priceupdater/stocks-api/middleware/corel"
	"priceupdater/stocks-api/utils/logging"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

func provideInstrumentRouter() *instrumentRouter {
	dbrepo := db.NewInstrumentRepo()
	numberOfWorkerStr := os.Getenv("SERVICE_WORKER_POOL_SIZE")
	numberOfWorker, _ := strconv.Atoi(numberOfWorkerStr)
	if numberOfWorker == 0 {
		numberOfWorker = 50
	}
	instrumentService := usecase.NewInstrumentService(http.NewInstrumentHttpRepo(), websocket.NewWebsocketRepo(dbrepo), dbrepo, worker.NewWorkerPool(numberOfWorker, 2*numberOfWorker))
	ctx := corel.CreateNewContext()
	logging.Logger.WriteLogs(ctx, "starting_initial_update_equity_details_job", logging.InfoLevel, logging.Fields{})
	_, err := instrumentService.UpdateEquityStockDetails(ctx)
	if err != nil {
		logging.Logger.WriteLogs(ctx, "unable_to_run_initial_equity_update_job", logging.ErrorLevel, logging.Fields{"error": err})
		time.Sleep(5 * time.Second)
		log.Fatal(err)
	}
	go func() {
		logging.Logger.WriteLogs(ctx, "starting_initial_update_derivative_details_job", logging.InfoLevel, logging.Fields{})
		err := instrumentService.UpdateDerivativeStockDetails(ctx)
		if err != nil {
			logging.Logger.WriteLogs(ctx, "unable_to_run_initial_derivative_update_job", logging.ErrorLevel, logging.Fields{"error": err})
			time.Sleep(5 * time.Second)
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
