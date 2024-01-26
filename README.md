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

   `mutation CreatUser($email: String!, $password: String!) {
      createUser(email: $email, password: $password) {
        email,
        password,
        id,
      }
      
    }
  
  mutation CreateDeck {
    createDeck(input: {
      userId: "get uuid from create user",
      id:"",
      name: "create new quiz",
      lastUpdate: "10-24-2025",
      dateCreated: "10-22-1992",
      flashcards: [{front: "hello", back: "from the other side"}, {front: "guess who's back", back: "back again"}]
    }) {
      id,
      userId,
      lastUpdate,
      flashcards{
        front,
        back
      },
      name,
    }
  }
  
   query GetDeckByUser {
    getDeckByUser(userId: "get UUID from create user") {
      id,
      dateCreated,
      lastUpdate,
      flashcards {
        front, back
      },
      name,
      userId
    }
  }
  
  subscription Comment {
    comment {
      id, userId, comment
    }
  }
  
  query GetDeckById {
    getDeckById(id: "get deck id from create deck") {
      name,
      userId,
      id,
      flashcards {
        front, back
      }
    }
  }
mutation DeleteUser($userId: String!){
  deleteUser(userId: $userId)
}`

## Considerations

- TODO implement sort and filter feature of elasticsearch search queries. 
