package irepo

import (
	"context"
	"sensibull/stocks-api/business/entities/core"
	"sensibull/stocks-api/business/entities/dto"
)

type IInstrumentRepo interface {
	UpsertInstrument(ctx context.Context, instrument core.Instrument) error
	DeleteInstrument(ctx context.Context, instrument core.Instrument) error
	GetInstrument(ctx context.Context, token int64) (core.Instrument, error)
	GetTokensForSymbol(ctx context.Context, isymbol, itype string) (core.Tokens, error)
	SaveTokenForSymbol(ctx context.Context, isymbol, itype string, tokens core.Tokens) error
}

type IInstrumentHttpRepo interface {
	GetUnderLyingDerivatives(ctx context.Context, underLyingToken int64) ([]dto.Instrument, error)
	GetUnderLyings(ctx context.Context) ([]dto.Instrument, error)
}

type IWebsocketRepo interface {
	WSEventListener(ctx context.Context) error
	Subscribe(ctx context.Context) error
	UnSubscribe(ctx context.Context) error
}
