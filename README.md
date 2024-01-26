# What is this

This is the backend server for the **[flashy cards app](https://github.com/Step-henC/flashycards-ui/tree/master)**.
This server utilizes an elasticsearch DB and sends comment data to kafka.

# How to run Flashy Cards backend

First: have docker engine running.
command line in project root directory.
`cd kafka-docker` and then command line: `docker-compose up -d`
Check `localhost:9021` for confluence. Create an topic named `flash-deck` and another topic `flash-deck-comment'

Back in root directory (`cd ..`), command line: `go run server.go` and check browser for `localhost:8080`.

Below are GraphQL requests to help get started:
![Screenshot (13)](https://github.com/Step-henC/flashycards-backend/assets/98792412/421ddf8b-aa85-4d90-ab59-896611eee047)

![Screenshot (14)](https://github.com/Step-henC/flashycards-backend/assets/98792412/aa5cd79a-1160-4eea-b7d4-6aeca6d7a58a)


## Considerations

- TODO implement sort and filter feature of elasticsearch search queries. 
