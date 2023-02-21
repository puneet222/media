# mediamonks

## How to run the server

### Prerequisites
*Docker[https://www.docker.com/products/docker-desktop/] should be installed on the system*
*docker-compose[https://docs.docker.com/compose/install/other/]*

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
I'm using go sync.Mutex[https://pkg.go.dev/sync#Mutex] to handle multiple concurrent request to server for saving key value pair.

## Handling Server Restart to persist key value pairs
I'm using mongodb here for persistent storage of key value pairs, we can also use file system as well but just for learning purposes I've used mongo.

## Speeding up the saving process of key value pair
So after saving the key value pair in the in-memory data store I'm firing up a new go routine that will save the same key value pair in mongodb making this process async will enhance the performance.

