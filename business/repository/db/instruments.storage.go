package db

import (
	"context"
	"sensibull/stocks-api/business/entities/core"
	"sensibull/stocks-api/business/interfaces/irepo"
	"sensibull/stocks-api/business/utility"
	"sensibull/stocks-api/db"
	"sensibull/stocks-api/utils/logging"
	"sync"
	"time"
)

type instrumentrepo struct {
	cache
}

var once sync.Once
var repo *instrumentrepo

func NewInstrumentRepo() irepo.IInstrumentRepo {
	once.Do(func() {
		repo = &instrumentrepo{}
		repo.redisCli = db.GetRedisClient()
	})
	return repo
}

func (ir *instrumentrepo) UpsertInstrument(ctx context.Context, instrument core.Instrument) error {
	dur := 2 * time.Minute
	if instrument.InstrumentType == "EQ" {
		dur = 20 * time.Minute
	}
	// update symbol -> token mapping in db here
	err := ir.encache(ctx, instrument.Symbol, instrument.Token, 0)
	if err != nil {
		logging.Logger.WriteLogs(ctx, "error_mapping_symbol_to_token", logging.ErrorLevel, logging.Fields{"error": err})
	}
	return ir.encache(ctx, utility.GetInstrumentKey(instrument.Token), instrument, dur)
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

func (ir *instrumentrepo) GetTokensAgainstToken(ctx context.Context, itoken string, itype string) (core.Tokens, error) {
	tokens := core.Tokens{}
	err := ir.read(ctx, utility.GetTokenKey(itoken, itype), &tokens)
	if err != nil {
		return tokens, err
	}
	return tokens, nil
}

func (ir *instrumentrepo) SaveTokensAgainstToken(ctx context.Context, itoken, itype string, tokens core.Tokens) error {
	err := ir.encache(ctx, utility.GetTokenKey(itoken, itype), tokens, 0)
	if err != nil {
		return err
	}
	return nil
}

func (ir *instrumentrepo) GetInstrumentToken(ctx context.Context, symbol string) (int64, error) {
	var token int64
	err := ir.read(ctx, symbol, &token)
	if err != nil {
		logging.Logger.WriteLogs(ctx, "error_fetching_token_from_symbol", logging.ErrorLevel, logging.Fields{"error": err})
		return 0, err
	}
	return token, nil
}
