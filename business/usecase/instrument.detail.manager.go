package usecase

import (
	"context"
	"sensibull/stocks-api/business/entities/dto"
	"sensibull/stocks-api/business/interfaces/irepo"
	"sensibull/stocks-api/business/interfaces/iusecase"
	"sync"
)

type instrumentservice struct {
	httpir    irepo.IInstrumentHttpRepo
	websocket irepo.IWebsocketRepo
	db        irepo.IInstrumentRepo
}

var once sync.Once
var service *instrumentservice

func NewInstrumentService(httpir irepo.IInstrumentHttpRepo, websocket irepo.IWebsocketRepo, db irepo.IInstrumentRepo) iusecase.IStocksInstrumentsService {
	once.Do(func() {
		service = &instrumentservice{
			httpir:    httpir,
			websocket: websocket,
			db:        db,
		}
	})
	return service
}

func (is *instrumentservice) FetchEquityStockDetails(ctx context.Context) ([]dto.Instrument, error) {
	return nil, nil
}

func (is *instrumentservice) FetchDerivativeStockDetails(ctx context.Context, symbol string) ([]dto.Instrument, error) {
	return nil, nil
}

func (is *instrumentservice) UpdateEquityStockDetails(ctx context.Context) error {
	return nil
}

func (is *instrumentservice) UpdateDerivativeStockDetails(ctx context.Context) error {
	return nil
}

func (is *instrumentservice) UpdateInstrumentPrice(ctx context.Context) error {
	return nil
}
