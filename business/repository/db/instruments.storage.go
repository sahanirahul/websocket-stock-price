package db

import (
	"context"
	"fmt"
	"sensibull/stocks-api/business/entities/core"
	"sensibull/stocks-api/business/interfaces/irepo"
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

func getInstrumentKey(instrument core.Instrument) string {
	// return fmt.Sprintf("%s:%d", "instrument_token:", instrument.Token)
	return fmt.Sprint(instrument.Token)
}

func (ir *instrumentrepo) UpsertInstrument(ctx context.Context, instrument core.Instrument) error {
	return ir.encache(ctx, getInstrumentKey(instrument), instrument, 0)
}

func (ir *instrumentrepo) DeleteInstrument(ctx context.Context, instrument core.Instrument) error {
	return ir.delete(ctx, getInstrumentKey(instrument))
}

func (ir *instrumentrepo) GetInstrument(ctx context.Context, instrument core.Instrument) (core.Instrument, error) {
	val := core.Instrument{}
	err := ir.read(ctx, getInstrumentKey(instrument), &val)
	if err != nil {
		return val, err
	}
	return val, nil
}
