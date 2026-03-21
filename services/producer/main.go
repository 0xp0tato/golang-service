package main

import (
	"log"
	"net/url"
	"os"
	"shared/kafka"

	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	wsurl := url.URL{Scheme: "ws", Host: os.Getenv("WEBSOCKET_URL"), Path: "/ws"}
	log.Printf("Connecting to %s", wsurl.String())

	// Connect to WebSocket server
	conn, _, err := websocket.DefaultDialer.Dial(wsurl.String(), nil)
	if err != nil {
		log.Fatal("Cannot connect to data-generator:", err)
	}
	defer conn.Close()

	kafkaurl := os.Getenv("KAFKA_URL")

	//Create KAFKA producer
	kafkaProducer, err := kafka.NewKafkaProducer(
		[]string{kafkaurl},
		"data-topic",
	)

	if err != nil {
		log.Fatal("Failed to create Kafka producer:", err)
	}
	defer kafkaProducer.CloseProducer()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("error reading message:", err)
			return
		}

		if err := kafkaProducer.SendData(message); err != nil {
			log.Printf("Failed to send message: %v", err)
		} else {
			log.Printf("Sent message to Kafka")
		}
	}
}
