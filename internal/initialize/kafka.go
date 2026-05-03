package initialize

import (
	"go-learning/global"

	kafka "github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

func InitKafka() {
	global.KafkaProducer = &kafka.Writer{
		Addr:     kafka.TCP("localhost:19092"),
		Topic:    "otp-auth-topic",
		Balancer: &kafka.LeastBytes{},
	}
}

func CloseKafka() {
	if err := global.KafkaProducer.Close(); err != nil {
		global.Logger.Error("Failed to close Kafka producer", zap.Error(err))
	}
}
