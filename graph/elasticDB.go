package graph

import (
	elasticsearch7 "github.com/elastic/go-elasticsearch/v7"
)

func NewClient() *elasticsearch7.Client { //must be capitlized to be visible by other packages

	goClient, err := elasticsearch7.NewClient(elasticsearch7.Config{
		Addresses: []string{"http://localhost:9200", "http://localhost:9201"},
		Username:  "",
		Password:  "",
	})

	if err != nil {

		panic(err)
	}

	goClient.Indices.Create("flash-deck") //create index on start up name of our kafka topic
	//create index here since struggled with index API of elasticsearhc
	//tried postman and curling in the docker container

	goClient.Indices.Create("flash-deck-deck")    //hold deck by users
	goClient.Indices.Create("flash-deck-comment") //hold user comments

	return goClient
}
