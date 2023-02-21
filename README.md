# mediamonks

## How to run the server

### Prerequisites
*[Docker](https://www.docker.com/products/docker-desktop/) and [docker compose](https://docs.docker.com/compose/install/other/) should be installed on the system*

### Follow the below steps

- Clone the repository
```
git clone https://github.com/puneet222/media.git
```
- Run docker compose command
```
docker compose up -d
```

## APIs

- **REST Endpoint to fetch the value of a given key**
```
curl 'http://localhost:8088/data?key=key1'
```
*Key will be passed as query params to the url*

- **Websocket API**

```
http://localhost:8088/ws
```
*URL for websocket connection*

```
{
    "key": "key1",
    "value": "value1"
}
```
*Request Message to be sent over to websocket*
```
{
    "success": true
    "message": string
}
```
*Response message received from websocket*

## Handling concurrency
To handle multiple concurrent requests to the server for saving key-value pairs, I'm using the [sync.Mutex](https://pkg.go.dev/sync#Mutex) package in Go.

## Persisting Key-Value Pairs Across Server Restarts
I'm utilizing MongoDB for persistent storage of key-value pairs to ensure they are saved across server restarts.

## Speeding Up the Key-Value Pair Saving Process
To improve performance, I'm using a new Go routine to save the same key-value pair in MongoDB asynchronously after saving it in the in-memory data store:
```
go app.InsertToDB(keyValue)
```

### Additional thoughts

- Instead of utilizing an in-memory data store (e.g. a map), we could consider using Redis to store key-value pairs.
- While I have been occupied with work, I am currently conducting tests on the database (DB) and will ensure the repository is complete by adding the remaining tests.
- For our assignment, we could also opt to use the file system instead of MongoDB for persistent storage.
- As for testing the REST API and WebSocket API, I utilized Postman.