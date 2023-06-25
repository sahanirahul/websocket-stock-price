package connection

import (
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

var webSocketConn *websocket.Conn

var once sync.Once

func GetWebSocketConnection(websocketURL string, establishNewConn bool) (*websocket.Conn, error) {
	once.Do(func() {
		conn, _, err := websocket.DefaultDialer.Dial(websocketURL, nil)
		if err != nil {
			log.Fatal("WebSocket connection error:", err)
		}
		webSocketConn = conn
	})
	if establishNewConn {
		conn, _, err := websocket.DefaultDialer.Dial(websocketURL, nil)
		if err != nil {
			log.Fatal("WebSocket connection error:", err)
		}
		webSocketConn = conn
	}
	return webSocketConn, nil
}
