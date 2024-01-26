package consumer

import (
	"context"
	"encoding/json"
	"flashy-cards-kafka-producer/graph/model"
	"fmt"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/segmentio/kafka-go"
)

const (
	brokerAddress = "127.0.0.1:9092"
)

func Consume(ctx context.Context, elasticSearchClient *elasticsearch.Client) model.Comment {
	// initialize a new reader with the brokers and topic
	// the groupID identifies the consumer and prevents
	// it from receiving duplicate messages
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     []string{brokerAddress},
		Topic:       "flash-deck-comment",
		GroupID:     "only-flash-group",
		StartOffset: kafka.LastOffset, //new messages
	})
	for {
		// the `ReadMessage` method blocks until we receive the next event
		msg, err := r.ReadMessage(ctx)
		if err != nil {
			panic("could not read message " + err.Error())
		}
		// after receiving the message, log its value
		fmt.Println("received: ", string(msg.Value))

		if err != nil {
			fmt.Println(err)
		}

		//elasticSearchClient.Index("flash-deck-comment", bytes.NewReader(msg.Value), elasticSearchClient.Index.WithDocumentID(string(msg.Key)))

		var comm model.Comment
		err = json.Unmarshal(msg.Value, &comm)
		if err != nil {
			fmt.Println(err)
		}
		return comm
	}

}
