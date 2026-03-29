package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"os"
	"shared/db"
	database "shared/db"
	"shared/enum"
	"shared/kafka"
	"shared/metrics"

	util "shared/utils"

	"github.com/IBM/sarama"
)

func consume(partitionConsumer sarama.PartitionConsumer, sqldb *sql.DB) error {
	metrics.StartMetricsServer("2114")
	var msgConsumed = metrics.NewCounter("consumer_messages_consumed_total", "Total messages consumed")
	var msgFailed = metrics.NewCounter("consumer_messages_failed_total", "Total messages failed")

	signals := util.GetSignalChannel()
	defer sqldb.Close()
	for {
		select {
		case msg := <-partitionConsumer.Messages():
			var data enum.Data
			if err := json.Unmarshal(msg.Value, &data); err != nil {
				log.Printf("Error unmarshaling message: %v", err)
				msgFailed.Inc()
				continue
			}

			// Process your data here
			db.Insert(sqldb, data)
			msgConsumed.Inc()

		case err := <-partitionConsumer.Errors():
			log.Printf("Error: %v", err)
			msgFailed.Inc()

		case <-signals:
			log.Println("Shutting down consumer...")
			return nil
		}
	}
}

func main() {
	db := database.ConnectDB()
	defer database.DeleteTable(db)

	kafkaurl := os.Getenv("KAFKA_URL")

	kafkaConsumer, err := kafka.NewKafkaConsumer(
		[]string{kafkaurl},
		"data-topic",
	)

	if err != nil {
		log.Fatal("Failed to create Kafka consumer:", err)
	}
	defer kafkaConsumer.CloseConsumer()

	log.Println("Starting Kafka consumer...")
	partitionConsumer, err := kafkaConsumer.ConsumeMessages()
	if err != nil {
		log.Fatal("Error consuming messages:", err)
	}
	defer partitionConsumer.Close()

	consume(partitionConsumer, db)
}
