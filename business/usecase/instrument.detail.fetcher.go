package usecase

import (
	"context"
	"fmt"
	"priceupdater/stocks-api/business/entities/core"
	"priceupdater/stocks-api/business/entities/dto"
	"priceupdater/stocks-api/business/interfaces/icore"
	"priceupdater/stocks-api/business/interfaces/irepo"
	"priceupdater/stocks-api/business/interfaces/iusecase"
	"priceupdater/stocks-api/business/utility"
	"priceupdater/stocks-api/utils/logging"
	"sync"
)

type instrumentservice struct {
	httpir    irepo.IInstrumentHttpRepo
	websocket irepo.IWebsocketRepo
	db        irepo.IInstrumentRepo
	pool      icore.IPool
}

var once sync.Once
var service *instrumentservice

func NewInstrumentService(httpir irepo.IInstrumentHttpRepo, websocket irepo.IWebsocketRepo, db irepo.IInstrumentRepo, pool icore.IPool) iusecase.IStocksInstrumentsService {
	once.Do(func() {
		service = &instrumentservice{
			httpir:    httpir,
			websocket: websocket,
			db:        db,
			pool:      pool,
		}
	})
	return service
}

func (is *instrumentservice) FetchEquityStockDetails(ctx context.Context) ([]dto.Instrument, error) {
	tokens, err := is.db.GetTokensAgainstToken(ctx, utility.TOKENFORALLUNDERLYING, utility.EQUITY)
	if err != nil {
		return nil, err
	}
	if tokens.Set == nil {
		return nil, nil
	}
	instruments := []core.Instrument{}
	for _, token := range tokens.Set.Values() {
		ins, err := is.db.GetInstrument(ctx, int64(token.(float64)))
		if err != nil {
			logging.Logger.WriteLogs(ctx, "error_fetching_equity_instrument_from_token", logging.ErrorLevel, logging.Fields{"error": err})
			continue
		}
		if len(ins.Symbol) > 0 && ins.Price > 0 {
			instruments = append(instruments, ins)
		}
	}
	return core.GetDtoInstruments(instruments), nil
}

func (is *instrumentservice) FetchDerivativeStockDetails(ctx context.Context, symbol string) ([]dto.Instrument, error) {
	// todo: fetch the token for the symbol
	underlyingToken, err := is.db.GetInstrumentToken(ctx, symbol)
	if err != nil {
		return nil, fmt.Errorf("error fetching derivatives for %s", symbol)
	}
	derivativeTokens, err := is.db.GetTokensAgainstToken(ctx, fmt.Sprint(underlyingToken), utility.DERIVATIVES)
	if err != nil {
		return nil, err
	}
	if derivativeTokens.Set == nil {
		return nil, nil
	}
	instruments := []core.Instrument{}
	for _, token := range derivativeTokens.Set.Values() {
		ins, err := is.db.GetInstrument(ctx, int64(token.(float64)))
		if err != nil {
			logging.Logger.WriteLogs(ctx, "error_fetching_derivative_instrumnet_from_token", logging.ErrorLevel, logging.Fields{"error": err})
			continue
		}
		if len(ins.Symbol) > 0 && ins.Price > 0 {
			instruments = append(instruments, ins)
		}
	}
	return core.GetDtoInstruments(instruments), nil
}
