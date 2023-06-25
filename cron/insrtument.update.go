package cron

import (
	"context"
	"sensibull/stocks-api/business/interfaces/iusecase"
	"sensibull/stocks-api/business/repository/db"
	"sensibull/stocks-api/business/repository/http"
	"sensibull/stocks-api/business/repository/websocket"
	"sensibull/stocks-api/business/usecase"
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
		cronObj.instrumentService = usecase.NewInstrumentService(http.NewInstrumentHttpRepo(), websocket.NewWebsocketRepo(), db.NewInstrumentRepo())
	})
	cronObj.startUnderlyingUpdate()
	time.Sleep(10 * time.Second)
	cronObj.startUnderlyingDerivativeUpdate()
	return cronObj
}

func (cro *cronn) startUnderlyingUpdate() {
	c := cron.New()
	c.AddFunc("*/15 * * * *", func() {
		ctx := context.Background()
		logging.Logger.WriteLogs(ctx, "cron_started_equity", logging.InfoLevel, logging.Fields{})
		err := cro.instrumentService.UpdateEquityStockDetails(ctx)
		if err != nil {
			logging.Logger.WriteLogs(ctx, "error_while_executing_equity_cron", logging.ErrorLevel, logging.Fields{"error": err})
		}
		logging.Logger.WriteLogs(ctx, "cron_finished_equity", logging.InfoLevel, logging.Fields{})
	})
	c.Start()
}

func (cro *cronn) startUnderlyingDerivativeUpdate() {
	c := cron.New()
	c.AddFunc("*/1 * * * *", func() {
		ctx := context.Background()
		logging.Logger.WriteLogs(ctx, "cron_started_derivatives", logging.InfoLevel, logging.Fields{})
		err := cro.instrumentService.UpdateDerivativeStockDetails(ctx)
		if err != nil {
			logging.Logger.WriteLogs(ctx, "error_while_executing_equity_cron", logging.ErrorLevel, logging.Fields{"error": err})
		}
		logging.Logger.WriteLogs(ctx, "cron_finished_derivatives", logging.InfoLevel, logging.Fields{})
	})
	c.Start()
}
