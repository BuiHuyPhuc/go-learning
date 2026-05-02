package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	kafka "github.com/segmentio/kafka-go"
)

var (
	kafkaProducer *kafka.Writer
)

const (
	kafkaURL   = "localhost:19092"
	kafkaTopic = "user_topic_001"
)

// for producer
func getKafkaWriter(kafkaURL, topic string) *kafka.Writer {
	return &kafka.Writer{
		Addr:     kafka.TCP(kafkaURL),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
}

// for consumer
func getKafkaReader(kafkaURL, topic, groupID string) *kafka.Reader {
	brokers := strings.Split(kafkaURL, ",")
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:        brokers, // []string{"localhost1:9092", "localhost2:9092", "localhost3:9092"},
		GroupID:        groupID,
		Topic:          topic,
		MinBytes:       10e3, // 10KB
		MaxBytes:       10e6, // 10MB
		CommitInterval: time.Second,
		StartOffset:    kafka.LastOffset, // FirstOffset
	})
}

type StockInfo struct {
	Message string `josn:"message"`
	Type    string `json:"type"`
}

func newStock(msg, typeMsg string) *StockInfo {
	s := StockInfo{}
	s.Message = msg
	s.Type = typeMsg

	return &s
}

func actionStock(c *gin.Context) {
	s := newStock(c.Query("msg"), c.Query("type"))
	body := make(map[string]interface{})
	body["action"] = "action"
	body["info"] = s

	jsonBody, _ := json.Marshal(body)

	// tạo msg
	msg := kafka.Message{
		Key:   []byte("action"),
		Value: []byte(jsonBody),
	}

	// viết msg vào kafka
	err := kafkaProducer.WriteMessages(context.Background(), msg)
	if err != nil {
		c.JSON(200, gin.H{
			"err": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"err": "",
		"msg": "action successfully",
	})
}

func RegisterConsumerATC(id int) {
	// group consumer
	kafkaGroupId := fmt.Sprintf("consumer-%d", id)
	reader := getKafkaReader(kafkaURL, kafkaTopic, kafkaGroupId)
	defer reader.Close()

	fmt.Printf("Consumer (%d) Hong Phien ATC::\n", id)
	for {
		m, err := reader.ReadMessage(context.Background())
		if err != nil {
			fmt.Printf("Consumer (%d) error: %v", id, err)
			return
		}
		fmt.Printf("Consumer (%d), hong topic: %v, partition: %v, offset: %v, time: %d %s = %s\n",
			id, m.Topic, m.Partition, m.Offset, m.Time.Unix(), string(m.Key), string(m.Value))
	}
}

func main() {
	r := gin.Default()
	kafkaProducer = getKafkaWriter(kafkaURL, kafkaTopic)
	defer kafkaProducer.Close()

	r.POST("action/stock", actionStock)

	// đăng ký 2 user để mua Stock trong ATC (1) (2)
	go RegisterConsumerATC(1)
	go RegisterConsumerATC(2)
	go RegisterConsumerATC(3)

	r.Run(":8999")
}
