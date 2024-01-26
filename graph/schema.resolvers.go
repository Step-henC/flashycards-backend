package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.43

import (
	"bytes"
	"context"
	"encoding/json"
	"flashy-cards-kafka-producer/graph/model"
	"fmt"
	"strings"

	"github.com/google/uuid"
	kafka "github.com/segmentio/kafka-go"
)

// Elasticsearch automatically Upserts by id. Better name is upsert deck
func (r *mutationResolver) CreateDeck(ctx context.Context, input *model.NewDeck) (*model.Deck, error) {
	//add the cards for a FlashCards field
	var answerCards []*model.Cards
	for _, card := range input.Flashcards {
		cardToInsert := model.Cards{
			Front: card.Front,
			Back:  card.Back,
		}

		answerCards = append(answerCards, &cardToInsert)
	}

	//add flashcards and input to a Deck struct
	answerDeck := model.Deck{
		Name:        input.Name,
		UserID:      input.UserID,
		DateCreated: input.DateCreated,
		LastUpdate:  input.LastUpdate,
		Flashcards:  answerCards,
	}
	//set id if new, or upsert deck with current id if it exists
	if input.ID == "" {
		answerDeck.ID = uuid.NewString()
	} else {
		answerDeck.ID = input.ID
	}

	//marshal deck to send to elasticsearch via http
	marshaledDeck, prob := json.Marshal(answerDeck)
	if prob != nil {
		return &model.Deck{}, prob
	}
	//save deck to elasticsearch
	_, err := r.DB.Index("flash-deck-deck", bytes.NewReader(marshaledDeck), r.DB.Index.WithDocumentID(answerDeck.ID))
	if err != nil {
		return &model.Deck{}, err
	}

	//return deck
	return &answerDeck, nil
}

// originally was a kafka topic, but for now just save user to DB
func (r *mutationResolver) CreateUser(ctx context.Context, email string, password string) (*model.User, error) {
	//prod := producer.NewProducer()

	user := model.User{
		ID:       uuid.NewString(),
		Email:    email,
		Password: password,
	}

	data, prob := json.Marshal(user)

	if prob != nil {
		return &model.User{}, prob
	}

	//save deck to elasticsearch
	_, err := r.DB.Index("flash-deck", bytes.NewReader(data), r.DB.Index.WithDocumentID(user.ID))
	if err != nil {
		return &model.User{}, err
	}
	//prod.CreateTopic("flash-deck")

	//ok := prod.Produce([]byte(user.ID), data, "flash-deck") //use a unique id for key! uuid here works

	// if ok != nil {
	// 	return &model.User{}, ok
	// }

	// res, err := r.DB.Index("flash-deck", bytes.NewReader(data), r.DB.Index.WithDocumentID(user.ID))

	// if err != nil {
	// 	return &model.User{}, err
	// }

	// defer res.Body.Close()

	return &user, nil
}

// DeleteDeckByUser is the resolver for the deleteDeckByUser field.
func (r *mutationResolver) DeleteDeckByUser(ctx context.Context, deckID string) (string, error) {
	query := `{ "query": { "match": {"id.keyword": "` + deckID + `"} } }`

	//stringfy query
	var builder strings.Builder
	builder.WriteString(query)
	_, err := r.DB.DeleteByQuery([]string{"flash-deck-deck"}, strings.NewReader(builder.String()))

	if err != nil {
		return "", err
	}

	return deckID, nil
}

// DeleteUser is the resolver for the deleteUser field.
func (r *mutationResolver) DeleteUser(ctx context.Context, userID string) (string, error) {
	query := `{ "query": { "match": {"id.keyword": "` + userID + `"} } }`

	//stringfy query
	var builder strings.Builder
	builder.WriteString(query)
	_, err := r.DB.DeleteByQuery([]string{"flash-deck"}, strings.NewReader(builder.String()))

	if err != nil {
		return "", err
	}

	return userID, nil
}

// DeleteAllDecksByUser is the resolver for the deleteAllDecksByUser field.
func (r *mutationResolver) DeleteAllDecksByUser(ctx context.Context, userID string) (string, error) {
	query := `{ "query": { "match": {"userId.keyword": "` + userID + `"} } }`

	//stringfy query
	var builder strings.Builder
	builder.WriteString(query)
	_, err := r.DB.DeleteByQuery([]string{"flash-deck-deck"}, strings.NewReader(builder.String()))

	if err != nil {
		return "", err
	}

	return userID, nil
}

// GetDeckByUser is the resolver for the getDeckByUser field.
func (r *queryResolver) GetDeckByUser(ctx context.Context, userID string) ([]*model.Deck, error) {
	var answerDeck []*model.Deck

	query := `{ "query": { "match": {"userId.keyword": "` + userID + `"} } }`

	//stringfy query
	var builder strings.Builder
	builder.WriteString(query)

	results, err := r.DB.Search(r.DB.Search.WithIndex("flash-deck-deck"),
		r.DB.Search.WithContext(context.Background()),
		r.DB.Search.WithTrackTotalHits(true),
		r.DB.Search.WithBody(strings.NewReader(builder.String())))

	if err != nil {
		return []*model.Deck{}, err
	}

	var mapResp map[string]interface{}

	problem := json.NewDecoder(results.Body).Decode(&mapResp)
	if problem != nil {
		return []*model.Deck{}, problem
	}

	defer results.Body.Close()

	data := mapResp["hits"].(map[string]any)

	var deckList []interface{}

	for i := 0; i < len(data); i++ {

		if i == 2 {

			for _, person := range data["hits"].([]interface{}) {

				deckList = append(deckList, person.(map[string]interface{})["_source"])

			}

		}

	}

	for _, cards := range deckList {

		res := cards.(map[string]interface{})

		if res["userId"].(string) == userID {

			var flashyCards []*model.Cards
			for _, v := range res["flashcards"].([]interface{}) {

				defineV := v.(map[string]interface{})

				g := model.Cards{
					Front: defineV["front"].(string),
					Back:  defineV["back"].(string),
				}

				flashyCards = append(flashyCards, &g)
			}
			deckToCollect := model.Deck{
				UserID:      userID,
				Name:        res["name"].(string),
				ID:          res["id"].(string),
				Flashcards:  flashyCards,
				LastUpdate:  res["lastUpdate"].(string),
				DateCreated: res["dateCreated"].(string),
			}

			answerDeck = append(answerDeck, &deckToCollect)

		}

	}

	// null proof the data
	if len(answerDeck) == 0 {
		return []*model.Deck{}, nil
	}
	return answerDeck, nil
}

// GetSortedDeck is the resolver for the getSortedDeck field.
func (r *queryResolver) GetSortedDeck(ctx context.Context, options []string) ([]*model.Deck, error) {
	query := `{ "query": { "match": {"userId.keyword": "` + options[0] + `"} } }`

	//stringfy query
	var builder strings.Builder
	builder.WriteString(query)

	results, err := r.DB.Search(r.DB.Search.WithIndex("flash-deck-deck"),
		r.DB.Search.WithContext(context.Background()),
		r.DB.Search.WithTrackTotalHits(true),
		r.DB.Search.WithSort("lastUpdate", "asc"), //https://53jk1.medium.com/using-elasticsearch-in-golang-api-writing-a-searchquery-structure-for-http-requests-b6c99603aaf1
		r.DB.Search.WithBody(strings.NewReader(builder.String())))

	if err != nil {
		return []*model.Deck{}, err
	}

	var mapResp map[string]interface{}
	problem := json.NewDecoder(results.Body).Decode(&mapResp)
	if problem != nil {
		return []*model.Deck{}, problem
	}

	fmt.Println(mapResp)
	defer results.Body.Close()
	return []*model.Deck{}, nil
}

// GetUsers is the resolver for the getUsers field.
func (r *queryResolver) GetUsers(ctx context.Context, email string, password string) (string, error) {
	panic(fmt.Errorf("not implemented: GetUsers - getUsers"))
}

// GetDeckByID is the resolver for the getDeckById field.
func (r *queryResolver) GetDeckByID(ctx context.Context, id string) (*model.Deck, error) {
	answerDeck := model.Deck{Flashcards: []*model.Cards{}}

	query := `{ "query": { "match": {"id.keyword": "` + id + `"} } }`

	//stringfy query
	var builder strings.Builder
	builder.WriteString(query)

	results, err := r.DB.Search(r.DB.Search.WithIndex("flash-deck-deck"),
		r.DB.Search.WithContext(context.Background()),
		r.DB.Search.WithTrackTotalHits(true),
		r.DB.Search.WithBody(strings.NewReader(builder.String())))

	if err != nil {
		return &answerDeck, err
	}

	var mapResp map[string]interface{}

	problem := json.NewDecoder(results.Body).Decode(&mapResp)
	if problem != nil {
		return &answerDeck, problem
	}

	defer results.Body.Close()

	data := mapResp["hits"].(map[string]any)

	var deckList []interface{}

	for i := 0; i < len(data); i++ {

		if i == 2 {

			for _, person := range data["hits"].([]interface{}) {

				deckList = append(deckList, person.(map[string]interface{})["_source"])

			}

		}

	}

	for _, cards := range deckList {

		res := cards.(map[string]interface{})

		if res["id"].(string) == id {

			var flashyCards []*model.Cards
			for _, v := range res["flashcards"].([]interface{}) {

				defineV := v.(map[string]interface{})

				g := model.Cards{
					Front: defineV["front"].(string),
					Back:  defineV["back"].(string),
				}

				flashyCards = append(flashyCards, &g)
			}
			deckToCollect := model.Deck{
				UserID:      res["userId"].(string),
				Name:        res["name"].(string),
				ID:          id,
				Flashcards:  flashyCards,
				LastUpdate:  res["lastUpdate"].(string),
				DateCreated: res["dateCreated"].(string),
			}
			// return this deck, should be one
			return &deckToCollect, nil

		}

	}

	return &answerDeck, nil
}

// kafka consumer for our "comment" topic
func (r *subscriptionResolver) Comment(ctx context.Context) (<-chan *model.Comment, error) {
	re := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     []string{"127.0.0.1:9092"},
		Topic:       "flash-deck-comment",
		GroupID:     "only-flash-group",
		StartOffset: kafka.LastOffset, //new messages
	})

	ch := make(chan *model.Comment)

	go func() {
		defer close(ch)

		for {
			msg, err := re.ReadMessage(ctx)
			if err != nil {
				//panic("could not read message " + err.Error())
				fmt.Print("could not read message: " + err.Error())

			}
			// after receiving the message, log its value
			fmt.Println("received: ", string(msg.Value))

			if err != nil {
				fmt.Println(err)
			}

			//elasticSearchClient.Index("flash-deck-comment", bytes.NewReader(msg.Value), elasticSearchClient.Index.WithDocumentID(string(msg.Key)))

			var comm *model.Comment
			err = json.Unmarshal(msg.Value, &comm)
			if err != nil {
				fmt.Println(err)
			}

			select {
			case <-ctx.Done():
				fmt.Println("Subscription closed")
				return //end routine
			case ch <- comm:
				//our message went through, do nothing
			}
		}

	}()
	return ch, nil
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

// Subscription returns SubscriptionResolver implementation.
func (r *Resolver) Subscription() SubscriptionResolver { return &subscriptionResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }
