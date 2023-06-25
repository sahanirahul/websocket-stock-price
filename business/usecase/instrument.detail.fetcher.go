package usecase

import (
	"context"
	"fmt"
	"sensibull/stocks-api/business/entities/core"
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

const (
	TOKENFORALLUNDERLYING = "ALLUNDERYINGS"
	EQUITY                = "EQ"
	DERIVATIVES           = "DERIVATIVES"
)

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
	tokens, err := is.db.GetTokensAgainstToken(ctx, TOKENFORALLUNDERLYING, EQUITY)
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
			return nil, nil
		}
		instruments = append(instruments, ins)
	}
	return core.GetDtoInstruments(instruments), nil
}

func (is *instrumentservice) FetchDerivativeStockDetails(ctx context.Context, symbol string) ([]dto.Instrument, error) {
	// todo: fetch the token for the symbol
	underlyongToken, err := is.db.GetInstrumentToken(ctx, symbol)
	if err != nil {
		return nil, fmt.Errorf("error fetching derivatives for %s", symbol)
	}
	derivativeTokens, err := is.db.GetTokensAgainstToken(ctx, fmt.Sprint(underlyongToken), DERIVATIVES)
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
			return nil, nil
		}
		instruments = append(instruments, ins)
	}
	return core.GetDtoInstruments(instruments), nil
}
