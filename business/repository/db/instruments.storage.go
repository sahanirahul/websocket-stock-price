package db

import (
	"context"
	"fmt"
	"priceupdater/stocks-api/business/entities/core"
	"priceupdater/stocks-api/business/interfaces/irepo"
	"priceupdater/stocks-api/business/utility"
	"priceupdater/stocks-api/db"
	"priceupdater/stocks-api/utils/logging"
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

func (ir *instrumentrepo) InsertInstrument(ctx context.Context, instrument core.Instrument) error {
	dur := 2 * time.Minute
	if instrument.InstrumentType == utility.EQUITY {
		dur = 20 * time.Minute
	}
	// update symbol -> token mapping in db here
	err := ir.encache(ctx, instrument.Symbol, instrument.Token, dur, false)
	if err != nil {
		logging.Logger.WriteLogs(ctx, "error_mapping_symbol_to_token", logging.ErrorLevel, logging.Fields{"error": err})
	}
	return ir.encache(ctx, utility.GetInstrumentKey(instrument.Token), instrument, dur, false)
}

func (ir *instrumentrepo) DeleteInstrument(ctx context.Context, instrument core.Instrument) error {
	err := ir.delete(ctx, utility.GetInstrumentKey(instrument.Token))
	if err != nil {
		return err
	}
	err = ir.delete(ctx, utility.GetTokenKey(fmt.Sprint(instrument.Token), utility.DERIVATIVES))
	if err != nil {
		return err
	}
	err = ir.delete(ctx, instrument.Symbol)
	if err != nil {
		return err
	}
	return nil
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
	dur := 2 * time.Minute
	if itype == utility.EQUITY {
		dur = 20 * time.Minute
	}
	err := ir.encache(ctx, utility.GetTokenKey(itoken, itype), tokens, dur, false)
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

func (ir *instrumentrepo) UpdateInstrument(ctx context.Context, instrument core.Instrument) error {
	return ir.encache(ctx, utility.GetInstrumentKey(instrument.Token), instrument, 0, true)
}
