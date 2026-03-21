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

	util "shared/utils"

	"github.com/IBM/sarama"
	"github.com/joho/godotenv"
)

func consume(partitionConsumer sarama.PartitionConsumer, sqldb *sql.DB) error {
	signals := util.GetSignalChannel()
	defer sqldb.Close()
	for {
		select {
		case msg := <-partitionConsumer.Messages():
			var data enum.Data
			if err := json.Unmarshal(msg.Value, &data); err != nil {
				log.Printf("Error unmarshaling message: %v", err)
				continue
			}

			// Process your data here
			db.Insert(sqldb, data)

		case err := <-partitionConsumer.Errors():
			log.Printf("Error: %v", err)

		case <-signals:
			log.Println("Shutting down consumer...")
			return nil
		}
	}
}

func main() {

	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

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
