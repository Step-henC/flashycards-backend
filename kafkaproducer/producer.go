package producer

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/segmentio/kafka-go"
)

const (
	broker1Address = "localhost:9092"
)

type Producer struct {
	Writer *kafka.Writer
	Dialer *kafka.Dialer
}

func NewProducer() *Producer {
	dialer := &kafka.Dialer{
		Timeout:   10 * time.Second,
		DualStack: true,
	}
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{broker1Address},
		Dialer:  dialer,
	})

	return &Producer{
		Writer: writer,
	}
}

func (producer *Producer) Produce(key []byte, value []byte, topic string) error {

	err := producer.Writer.WriteMessages(context.TODO(), kafka.Message{
		Topic:  topic,
		Offset: 0,
		Key:    key,
		Value:  value,
	})

	if err != nil {
		fmt.Printf("Delivery failed %v", err)
		return err
	} else {
		fmt.Printf("Message delivered %s | key: %s\n ", topic, string(key))
		return nil

	}
}

func (producer *Producer) CreateTopic(kafkaTopic string) {
	conn, err := kafka.Dial("tcp", "localhost:9092")
	if err != nil {
		panic(err.Error())
	}
	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		panic(err.Error())
	}
	controllerConn, err := kafka.Dial("tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
	if err != nil {
		panic(err.Error())
	}
	defer controllerConn.Close()

	topicConfigs := []kafka.TopicConfig{{Topic: kafkaTopic, NumPartitions: 1, ReplicationFactor: 1}}

	err = controllerConn.CreateTopics(topicConfigs...)
	if err != nil {
		panic(err.Error())
	}
}
