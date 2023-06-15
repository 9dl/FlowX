package FlowX

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

const (
	readBufferSize  = 1024
	writeBufferSize = 1024
)

func CreateWebSocketConnection(url string) (*websocket.Conn, error) {
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	return conn, err
}

func CloseWebSocketConnection(conn *websocket.Conn, closeCode int, closeMessage string) error {
	return conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(closeCode, closeMessage))
}

func ReadWebSocketMessage(conn *websocket.Conn) (messageType int, p []byte, err error) {
	return conn.ReadMessage()
}

func WriteWebSocketMessage(conn *websocket.Conn, messageType int, data []byte) error {
	return conn.WriteMessage(messageType, data)
}

func WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Upgrade(w, r, nil, readBufferSize, writeBufferSize)
	if err != nil {
		log.Printf("Error upgrading to WebSocket connection: %v\n", err)
		return
	}
	defer conn.Close()

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading WebSocket message:", err)
			break
		}

		fmt.Println("Received message:", string(message))

		err = conn.WriteMessage(messageType, message)
		if err != nil {
			log.Println("Error writing WebSocket message:", err)
			break
		}
	}
}
