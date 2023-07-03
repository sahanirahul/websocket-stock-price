package websocket

import (
	"encoding/json"
	"log"
	"os"
	"priceupdater/stocks-api/business/entities/core"
	"priceupdater/stocks-api/business/interfaces/icore"
	"priceupdater/stocks-api/business/interfaces/irepo"
	"priceupdater/stocks-api/business/repository/websocket/connection"
	"priceupdater/stocks-api/business/utility"
	"priceupdater/stocks-api/business/worker"
	"priceupdater/stocks-api/middleware"
	"priceupdater/stocks-api/middleware/corel"
	"priceupdater/stocks-api/utils/logging"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
)

const websocketURL = "wss://prototype.sbulltech.com/api/ws"

type websocketrepo struct {
	db   irepo.IInstrumentRepo
	pool icore.IPool
}

var once sync.Once
var repo *websocketrepo

func NewWebsocketRepo(db irepo.IInstrumentRepo) irepo.IWebsocketRepo {
	once.Do(func() {
		numberOfWorkerStr := os.Getenv("WEBSOCKET_WORKER_POOL_SIZE")
		numberOfWorker, _ := strconv.Atoi(numberOfWorkerStr)
		if numberOfWorker == 0 {
			numberOfWorker = 50
		}
		subscriptionChan = make(chan core.WebsocketSubscription, 100)
		repo = &websocketrepo{db: db, pool: worker.NewWorkerPool(numberOfWorker, 2*numberOfWorker)}
		go repo.updateSubscription()
		go repo.wsEventListener()
	})
	return repo
}

func (wr *websocketrepo) wsEventListener() error {
	ctx := corel.CreateNewContext()
	// adding recovery for websocket listener go routine
	defer func() {
		if err := recover(); err != nil {
			middleware.Recover(ctx, err)
		}
	}()
	newConn := false

	for {
		// new context for each message
		ctx = corel.CreateNewContext()
		conn, _ := connection.GetWebSocketConnection(websocketURL, newConn)
		if conn == nil {
			logging.Logger.WriteLogs(ctx, "error_getting_websocket_connection_read", logging.ErrorLevel, logging.Fields{})
			// reset the connection
			val := connectionRetryCount.Add(1)
			if val > 5 {
				logging.Logger.WriteLogs(ctx, "max_connection_retry_limit_exceeded", logging.ErrorLevel, logging.Fields{})
				log.Fatal("not able to get websocket connection. shutting down server")
			}
			newConn = true
			time.Sleep(time.Millisecond * 500)
			continue
		}
		connectionRetryCount.Add(-1 * connectionRetryCount.Add(0))
		newConn = false
		// Read message from WebSocket connection
		_, message, err := conn.ReadMessage()
		if err != nil {
			logging.Logger.WriteLogs(ctx, "websocket_read_errror", logging.ErrorLevel, logging.Fields{"error": err})
			continue
		}
		// Handle the received message
		// fmt.Println("Received message:", string(message))
		var event core.WebsocketPriceEvent
		err = json.Unmarshal(message, &event)
		if err != nil {
			logging.Logger.WriteLogs(ctx, "webSocket_unmarshal_error", logging.ErrorLevel, logging.Fields{"error": err})
			continue
		}
		if event.DataType == utility.DataTypeQuote {
			// logging.Logger.WriteLogs(ctx, "price_quote_event", logging.ErrorLevel, logging.Fields{"payload": event.Payload})
			job := core.NewJob(func() {
				ins, err := wr.db.GetInstrument(ctx, event.Payload.Token)
				if err != nil {
					logging.Logger.WriteLogs(ctx, "error_fetching_instrument_detail_for_price_update", logging.ErrorLevel, logging.Fields{"error": err})
					return
				}
				ins.Price = event.Payload.Price
				err = wr.db.UpdateInstrument(ctx, ins)
				if err != nil {
					logging.Logger.WriteLogs(ctx, "error_updating_instrument_price", logging.ErrorLevel, logging.Fields{"error": err, "instrument": ins})
					return
				}
			})
			wr.pool.AddJob(job)
		}
	}
}

var subscriptionChan chan core.WebsocketSubscription
var connectionRetryCount atomic.Int32

func (wr *websocketrepo) updateSubscription() error {
	ctx := corel.CreateNewContext()
	// adding recovery for websocket update subscription
	defer func() {
		if err := recover(); err != nil {
			middleware.Recover(ctx, err)
		}
	}()
	for {
		req := <-subscriptionChan
		conn, _ := connection.GetWebSocketConnection(websocketURL, false)
		if conn == nil {
			logging.Logger.WriteLogs(ctx, "error_getting_websocket_connection_write", logging.ErrorLevel, logging.Fields{})
			time.Sleep(time.Second)
			continue
		}
		payload, err := json.Marshal(req)
		if err != nil {
			return err
		}
		err = conn.WriteMessage(websocket.TextMessage, payload)
		if err != nil {
			logging.Logger.WriteLogs(ctx, "websocket_write_error", logging.ErrorLevel, logging.Fields{"error": err})
			time.Sleep(time.Second)
			continue
		}
	}

}

func (wr *websocketrepo) AddSubscriptionRequest(req core.WebsocketSubscription) {
	subscriptionChan <- req
}
