package db

import (
	"context"
	"encoding"
	"encoding/json"
	"errors"
	"sensibull/stocks-api/utils/logging"
	"time"

	"github.com/redis/go-redis/v9"
)

type cache struct {
	redisCli *redis.Client
}

var ErrInvalid error = errors.New("invalid_request")

func (che *cache) encache(ctx context.Context, key string, obj interface{}, dur time.Duration) error {
	if obj == nil || len(key) == 0 {
		return ErrInvalid
	}
	if _, ok := obj.(encoding.BinaryMarshaler); !ok {
		obj, _ = json.Marshal(obj)
	}
	err := che.redisCli.Set(ctx, key, obj, dur).Err()
	if err != nil {
		logging.Logger.WriteLogs(ctx, "cache-write-failed", logging.ErrorLevel, logging.Fields{"error": err, "key": key, "data": obj})
		return err
	}
	return nil
}

func (che *cache) delete(ctx context.Context, key string) error {
	if len(key) == 0 {
		return ErrInvalid
	}
	err := che.redisCli.Del(ctx, key).Err()
	if err == nil || err == redis.Nil {
		return nil
	}
	return err
}

func (che *cache) read(ctx context.Context, key string, dest interface{}) error {
	if len(key) == 0 {
		return ErrInvalid
	}
	res := che.redisCli.Get(ctx, key)
	if res.Err() == redis.Nil {
		return nil
	}
	if res.Err() != nil {
		logging.Logger.WriteLogs(ctx, "cache-read-failed", logging.ErrorLevel, logging.Fields{"error": res.Err(), "key": key})
		return res.Err()
	}
	err := res.Scan(&dest)
	if err != nil {
		raw, _ := res.Bytes()
		if err = json.Unmarshal(raw, &dest); err != nil {
			logging.Logger.WriteLogs(ctx, "cache-read-failed-scan", logging.ErrorLevel, logging.Fields{"error": err, "key": key})
			return err
		}
	}
	return nil
}
