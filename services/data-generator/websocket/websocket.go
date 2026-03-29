package websocket

import (
	"encoding/json"
	"fmt"
	"net/http"

	"shared/metrics"
	util "shared/utils"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var msgProduced = metrics.NewCounter("data_generator_messages_produced_total", "Total messages produced")

func WsHandler(w http.ResponseWriter, r *http.Request) {
	// Upgrade the HTTP connection to a WebSocket connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error upgrading:", err)
		return
	}
	defer conn.Close()

	for {
		data, err := json.Marshal(util.CreateData())
		msgProduced.Inc()
		if err != nil {
			fmt.Println("Error marshaling data:", err)
			break
		}
		if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
			fmt.Println("Error writing message:", err)
			break
		}
	}
}
