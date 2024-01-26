package producer

import (
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde/jsonschema"
)

const (
	bootstrapServers = "localhost:9092"
	urlRegistry      = "schema-registry:8081"
)

func NewProd(value interface{}) error {
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": bootstrapServers})
	if err != nil {
		return err
	}

	fmt.Println("new confluent producer")

	client, prob := schemaregistry.NewClient(schemaregistry.NewConfig(urlRegistry))

	if prob != nil {
		return prob
	}

	serializer, err := jsonschema.NewSerializer(client, serde.ValueSerde, jsonschema.NewSerializerConfig())

	if err != nil {
		return err
	}

	topic := "flash-deck"

	payload, err := serializer.Serialize(topic, &value)

	if err != nil {
		return err
	}

	deliveryChan := make(chan kafka.Event)

	err = p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          payload,
		Headers:        []kafka.Header{{Key: "myTestHeader", Value: []byte("who needs header")}},
	}, deliveryChan)

	if err != nil {
		return err
	}

	e := <-deliveryChan
	m := e.(*kafka.Message)

	if m.TopicPartition.Error != nil {
		fmt.Printf("Delivery failed: %v\n", m.TopicPartition.Error)
	} else {
		fmt.Printf("Delivered message to topic %s [%d] at offset %v\n",
			*m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
	}

	close(deliveryChan)
	return nil
}
