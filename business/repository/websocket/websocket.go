package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sensibull/stocks-api/business/entities/core"
	"sensibull/stocks-api/business/interfaces/irepo"
	"sensibull/stocks-api/business/repository/websocket/connection"
	"sensibull/stocks-api/business/utility"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
)

const websocketURL = "wss://prototype.sbulltech.com/api/ws"

type websocketrepo struct {
	db irepo.IInstrumentRepo
}

var once sync.Once
var repo *websocketrepo

func NewWebsocketRepo(db irepo.IInstrumentRepo) irepo.IWebsocketRepo {
	once.Do(func() {
		subscriptionChan = make(chan core.WebsocketSubscription, 10)
		repo = &websocketrepo{db: db}
		go repo.updateSubscription(context.Background())
		go repo.wsEventListener(context.Background())
	})
	return repo
}

func (wr *websocketrepo) wsEventListener(ctx context.Context) error {
	for {
		conn, err := connection.GetWebSocketConnection(websocketURL, false)
		if err != nil {
			log.Fatal("unable to get websocket connection")
		}
		// Read message from WebSocket connection
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("WebSocket read error:", err)
			continue
		}
		// Handle the received message
		fmt.Println("Received message:", string(message))
		var event core.WebsocketPriceEvent
		err = json.Unmarshal(message, &event)
		if err != nil {
			log.Println("WebSocket unmarshal error:", err)
			continue
		}
		if event.DataType == utility.DataTypeQuote {
			ctx := context.Background()
			ins, err := wr.db.GetInstrument(ctx, event.Payload.Token)
			if err != nil {
				log.Println("error fetching instrument for price update:", err)
				continue
			}
			ins.Price = event.Payload.Price
			err = wr.db.UpdateInstrument(ctx, ins)
			if err != nil {
				log.Println("error updating instrument price:", err)
				continue
			}
		}
	}
}

var subscriptionChan chan core.WebsocketSubscription
var retryCount atomic.Int32

func (wr *websocketrepo) updateSubscription(ctx context.Context) error {
	conn, err := connection.GetWebSocketConnection(websocketURL, false)
	if err != nil {
		log.Fatal("unable to get websocket connection")
		return err
	}
	for {
		req := <-subscriptionChan
		payload, err := json.Marshal(req)
		if err != nil {
			return err
		}
		fmt.Println(string(payload))
		err = conn.WriteMessage(websocket.TextMessage, payload)
		if err != nil {
			count := retryCount.Add(1)
			if count > 3 || websocket.IsCloseError(err, websocket.CloseNormalClosure) {
				conn, err = connection.GetWebSocketConnection(websocketURL, true)
				if err != nil {
					log.Fatal("unable to get websocket connection in write")
					return err
				}
				retryCount.Add(-1 * count)
			} else {
				log.Println("WebSocket write error:", err)
				time.Sleep(time.Second)
			}
			continue
		}
	}

}

func (wr *websocketrepo) AddSubscriptionRequest(req core.WebsocketSubscription) {
	subscriptionChan <- req
}
