package kafka

import (
	"github.com/IBM/sarama"
)

type KafkaProducer struct {
	producer sarama.SyncProducer
	topic    string
}

type KafkaConsumer struct {
	consumer sarama.Consumer
	topic    string
}

func NewKafkaProducer(brokers []string, topic string) (*KafkaProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}

	return &KafkaProducer{
		producer: producer,
		topic:    topic,
	}, nil
}

func (kp *KafkaProducer) SendData(data []byte) error {

	msg := &sarama.ProducerMessage{
		Topic: kp.topic,
		Value: sarama.ByteEncoder(data),
	}

	_, _, err := kp.producer.SendMessage(msg)
	return err
}

func (kp *KafkaProducer) CloseProducer() error {
	return kp.producer.Close()
}

func NewKafkaConsumer(brokers []string, topic string) (*KafkaConsumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	consumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		return nil, err
	}

	return &KafkaConsumer{
		consumer: consumer,
		topic:    topic,
	}, nil
}

func (kc *KafkaConsumer) CloseConsumer() error {
	return kc.consumer.Close()
}

func (kc *KafkaConsumer) ConsumeMessages() (sarama.PartitionConsumer, error) {
	partitionConsumer, err := kc.consumer.ConsumePartition(kc.topic, 0, sarama.OffsetNewest)
	if err != nil {
		return nil, err
	}

	return partitionConsumer, nil
}
