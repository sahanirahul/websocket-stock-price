package iusecase

import (
	"context"
	"sensibull/stocks-api/business/entities/dto"
)

type IStocksInstrumentsService interface {
	FetchEquityStockDetails(ctx context.Context) ([]dto.Instrument, error)
	FetchDerivativeStockDetails(ctx context.Context, symbol string) ([]dto.Instrument, error)
	// the below functions will be used to update the latest listing and update websocket subscription
	UpdateEquityStockDetails(ctx context.Context) error
	UpdateDerivativeStockDetails(ctx context.Context) error
}
