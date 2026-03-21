package main

import (
	"fmt"
	"net/http"
	"shared/metrics"

	"data-generator/websocket"
)

func main() {
	metrics.StartMetricsServer("2112")

	http.HandleFunc("/ws", websocket.WsHandler)
	fmt.Println("WebSocket server started on :8080")
	err := http.ListenAndServe("0.0.0.0:8080", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}

}
