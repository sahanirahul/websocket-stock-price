package db

import (
	"context"
	"sensibull/stocks-api/business/entities/core"
	"sensibull/stocks-api/business/interfaces/irepo"
	"sensibull/stocks-api/business/utility"
	"sync"

	"github.com/redis/go-redis/v9"
)

type instrumentrepo struct {
	cache
}

var once sync.Once
var repo *instrumentrepo

func NewInstrumentRepo(redisCli *redis.Client) irepo.IInstrumentRepo {
	once.Do(func() {
		repo = &instrumentrepo{}
	})
	return repo
}

func (ir *instrumentrepo) UpsertInstrument(ctx context.Context, instrument core.Instrument) error {
	return ir.encache(ctx, utility.GetInstrumentKey(instrument.Token), instrument, 0)
}

func (ir *instrumentrepo) DeleteInstrument(ctx context.Context, instrument core.Instrument) error {
	return ir.delete(ctx, utility.GetInstrumentKey(instrument.Token))
}

func (ir *instrumentrepo) GetInstrument(ctx context.Context, token int64) (core.Instrument, error) {
	val := core.Instrument{}
	err := ir.read(ctx, utility.GetInstrumentKey(token), &val)
	if err != nil {
		return val, err
	}
	return val, nil
}

func (ir *instrumentrepo) GetTokensForSymbol(ctx context.Context, isymbol, itype string) (core.Tokens, error) {
	tokens := core.Tokens{}
	err := ir.read(ctx, utility.GetTokenKey(isymbol, itype), &tokens)
	if err != nil {
		return tokens, err
	}
	return tokens, nil
}

func (ir *instrumentrepo) SaveTokenForSymbol(ctx context.Context, isymbol, itype string, tokens core.Tokens) error {
	err := ir.encache(ctx, utility.GetTokenKey(isymbol, itype), tokens, 0)
	if err != nil {
		return err
	}
	return nil
}
