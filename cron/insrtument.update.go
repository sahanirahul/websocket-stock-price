package cron

import (
	"sensibull/stocks-api/business/interfaces/iusecase"
	"sensibull/stocks-api/business/repository/db"
	"sensibull/stocks-api/business/repository/http"
	"sensibull/stocks-api/business/repository/websocket"
	"sensibull/stocks-api/business/usecase"
	"sensibull/stocks-api/business/worker"
	"sensibull/stocks-api/middleware"
	"sensibull/stocks-api/middleware/corel"
	"sensibull/stocks-api/utils/logging"
	"sync"
	"time"

	"github.com/robfig/cron"
)

type cronn struct {
	instrumentService iusecase.IInstrumentDetailManager
}

var once sync.Once
var cronObj *cronn

func NewCron() *cronn {
	once.Do(func() {
		cronObj = &cronn{}
		dbrepo := db.NewInstrumentRepo()
		cronObj.instrumentService = usecase.NewInstrumentService(http.NewInstrumentHttpRepo(), websocket.NewWebsocketRepo(dbrepo), dbrepo, worker.NewWorkerPool(50, 50))
	})
	cronObj.startUnderlyingUpdate()
	cronObj.startUnderlyingDerivativeUpdate()
	return cronObj
}

func (cro *cronn) startUnderlyingUpdate() {
	c := cron.New()
	c.AddFunc("@every 15m", func() {
		ctx := corel.CreateNewContext()
		// adding recovery for cron job
		defer func() {
			if err := recover(); err != nil {
				middleware.Recover(ctx, err)
			}
		}()
		logging.Logger.WriteLogs(ctx, "cron_started_equity", logging.InfoLevel, logging.Fields{})
		count := 1
		for count < 3 {
			shouldRetry, err := cro.instrumentService.UpdateEquityStockDetails(ctx)
			if err != nil {
				logging.Logger.WriteLogs(ctx, "error_while_executing_equity_cron", logging.ErrorLevel, logging.Fields{"error": err})
			}
			if shouldRetry {
				time.Sleep(time.Second)
				count++
			} else {
				break
			}
		}

		logging.Logger.WriteLogs(ctx, "cron_finished_equity", logging.InfoLevel, logging.Fields{})
	})
	c.Start()
}

func (cro *cronn) startUnderlyingDerivativeUpdate() {
	c := cron.New()
	c.AddFunc("0 * * * * *", func() {
		ctx := corel.CreateNewContext()
		// adding recovery for cron jobs
		defer func() {
			if err := recover(); err != nil {
				middleware.Recover(ctx, err)
			}
		}()
		logging.Logger.WriteLogs(ctx, "cron_started_derivatives", logging.InfoLevel, logging.Fields{})
		err := cro.instrumentService.UpdateDerivativeStockDetails(ctx)
		if err != nil {
			logging.Logger.WriteLogs(ctx, "error_while_executing_equity_cron", logging.ErrorLevel, logging.Fields{"error": err})
		}
		logging.Logger.WriteLogs(ctx, "cron_finished_derivatives", logging.InfoLevel, logging.Fields{})
	})
	c.Start()
}
