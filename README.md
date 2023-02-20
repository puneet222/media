# mediamonks

## How to run the server

### Prerequisites
*Docker should be installed on the system*

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
