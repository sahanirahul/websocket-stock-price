package connection

import (
	"context"
	"priceupdater/stocks-api/utils/logging"
	"sync"

	"github.com/gorilla/websocket"
)

var webSocketConn *websocket.Conn

var once sync.Once

func GetWebSocketConnection(websocketURL string, establishNewConn bool) (*websocket.Conn, error) {
	once.Do(func() {
		conn, _, err := websocket.DefaultDialer.Dial(websocketURL, nil)
		if err != nil {
			logging.Logger.WriteLogs(context.Background(), "error_connecting_websocket_server", logging.ErrorLevel, logging.Fields{"error": err, "url": websocketURL})
		}
		webSocketConn = conn
	})
	if establishNewConn {
		conn, _, err := websocket.DefaultDialer.Dial(websocketURL, nil)
		if err != nil {
			logging.Logger.WriteLogs(context.Background(), "error_connecting_websocket_server", logging.ErrorLevel, logging.Fields{"error": err, "url": websocketURL})
		}
		webSocketConn = conn
	}
	return webSocketConn, nil
}
