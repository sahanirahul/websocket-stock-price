package connection

import (
	"log"

	"github.com/gorilla/websocket"
)

var webSocketConn *websocket.Conn

func createWebSocketConnection(endpointURL string) error {
	// Create a new WebSocket connection
	conn, _, err := websocket.DefaultDialer.Dial(endpointURL, nil)
	if err != nil {
		log.Fatal("WebSocket connection error:", err)
		return err
	}
	webSocketConn = conn
	return nil
}

func GetWebSocketConnection(endpointURL string, establishNewConn bool) (*websocket.Conn, error) {
	if webSocketConn == nil || establishNewConn {
		err := createWebSocketConnection(endpointURL)
		if err != nil {
			return nil, err
		}
	}
	return webSocketConn, nil
}
