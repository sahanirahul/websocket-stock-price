package db

import (
	"context"
	"sensibull/stocks-api/business/entities/core"
	"sensibull/stocks-api/business/interfaces/irepo"
	"sync"

	"github.com/go-redis/redis"
)

type instrumentrepo struct {
	redisCli *redis.Client
}

var once sync.Once
var repo *instrumentrepo

func NewInstrumentRepo(redisCli *redis.Client) irepo.IInstrumentRepo {
	once.Do(func() {
		repo = &instrumentrepo{
			redisCli: redisCli,
		}
	})
	return repo
}

func (ar *instrumentrepo) UpsertInstrument(ctx context.Context, instrument core.Instrument) error {

	return nil
}

func (ar *instrumentrepo) DeleteInstrument(ctx context.Context, instrument core.Instrument) error {
	return nil
}
