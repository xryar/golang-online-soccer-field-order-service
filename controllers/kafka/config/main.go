package kafka

import (
	"order-service/config"
	"order-service/controllers/kafka"
	kafkaPayment "order-service/controllers/kafka/payment"

	"golang.org/x/exp/slices"
)

type KafkaConsummer struct {
	consumer *ConsumerGroup
	kafka    kafka.IKafkaRegistry
}

type IKafka interface {
	Register()
}

func NewKafkaConsumer(consumer *ConsumerGroup, kafka kafka.IKafkaRegistry) IKafka {
	return &KafkaConsummer{
		consumer: consumer,
		kafka:    kafka,
	}
}

func (kc *KafkaConsummer) Register() {
	kc.paymentHandler()
}

func (kc *KafkaConsummer) paymentHandler() {
	if slices.Contains(config.Config.Kafka.Topics, kafkaPayment.PaymentTopic) {
		kc.consumer.RegisterHandler(kafkaPayment.PaymentTopic, kc.kafka.GetPayment().HandlePayment)
	}
}
