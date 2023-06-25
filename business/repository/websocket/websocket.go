package websocket

import (
	"context"
	"sensibull/stocks-api/business/interfaces/irepo"
	"sync"
)

type websocketrepo struct {
}

var once sync.Once
var repo *websocketrepo

func NewWebsocketRepo() irepo.IWebsocketRepo {
	once.Do(func() {
		repo = &websocketrepo{}
	})
	return repo
}

func (cr *websocketrepo) WSEventListener(ctx context.Context) error {
	return nil
}

func (cr *websocketrepo) Subscribe(ctx context.Context) error {
	return nil
}
func (cr *websocketrepo) UnSubscribe(ctx context.Context) error {
	return nil
}
