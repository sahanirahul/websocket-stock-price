package irepo

import (
	"context"
	"sensibull/stocks-api/business/entities/core"
)

type IInstrumentRepo interface {
	InsertInstrument(ctx context.Context, instrument core.Instrument) error
	UpdateInstrument(ctx context.Context, instrument core.Instrument) error
	DeleteInstrument(ctx context.Context, instrument core.Instrument) error
	GetInstrument(ctx context.Context, token int64) (core.Instrument, error)
	GetInstrumentToken(ctx context.Context, symbol string) (int64, error)
	GetTokensAgainstToken(ctx context.Context, itoken, itype string) (core.Tokens, error)
	SaveTokensAgainstToken(ctx context.Context, itoken, itype string, tokens core.Tokens) error
}

type IInstrumentHttpRepo interface {
	GetUnderLyingDerivatives(ctx context.Context, underLyingToken int64) ([]core.Instrument, error)
	GetUnderLyings(ctx context.Context) ([]core.Instrument, error)
}

type IWebsocketRepo interface {
	AddSubscriptionRequest(req core.WebsocketSubscription)
}
