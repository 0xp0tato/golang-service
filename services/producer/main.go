package main

import (
	"log"
	"net/url"
	"os"
	"shared/kafka"
	"shared/metrics"
	"time"

	"github.com/gorilla/websocket"
)

func connectWebSocket(wsurl string) *websocket.Conn {
	for {
		conn, _, err := websocket.DefaultDialer.Dial(wsurl, nil)
		if err != nil {
			log.Println("Cannot connect to data-generator, retrying in 3s:", err)
			time.Sleep(3 * time.Second)
			continue
		}
		log.Println("Connected to WebSocket!")
		return conn
	}
}

func connectKafka(kafkaurl string) *kafka.KafkaProducer {
	for {
		kafkaProducer, err := kafka.NewKafkaProducer(
			[]string{kafkaurl},
			"data-topic",
		)
		if err != nil {
			log.Println("Failed to create Kafka producer, retrying in 3s:", err)
			time.Sleep(3 * time.Second)
			continue
		}
		log.Println("Connected to Kafka!")
		return kafkaProducer
	}
}

func main() {
	wsurl := url.URL{Scheme: "ws", Host: os.Getenv("WEBSOCKET_URL"), Path: "/ws"}
	log.Printf("Connecting to %s", wsurl.String())
	conn := connectWebSocket(wsurl.String())
	defer conn.Close()

	kafkaProducer := connectKafka(os.Getenv("KAFKA_URL"))
	defer kafkaProducer.CloseProducer()

	metrics.StartMetricsServer("2113")
	var msgPassed = metrics.NewCounter("producer_messages_passed_total", "Total messages passed")
	var msgFailed = metrics.NewCounter("producer_messages_passed_failed", "Total messages failed")

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("error reading message:", err)
			return
		}
		if err := kafkaProducer.SendData(message); err != nil {
			log.Printf("Failed to send message: %v", err)
			msgFailed.Inc()
		} else {
			log.Printf("Sent message to Kafka")
			msgPassed.Inc()
		}
	}
}
