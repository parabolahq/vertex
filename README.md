# Vertex

Message Queue pool for websocket messaging

## Getting Started

```bash
git clone https://github.com/parabolahq/vertex
```

### Prerequisites

* [Go 1.18+](https://go.dev/dl/)
* [RabbitMQ](https://www.rabbitmq.com/)

## Usage

1. Setup environment variables, or config via file
2. Sync dependencies
    ```bash
    go mod download && go mod verify
    ```
3. Build executable
    ```bash
    go build main.go
   ```
4. Run executable
    ```bash
   ./main
   ```

## Communication with service

### Errors

|       Int Error Code       | Description                             | Fix recommendations                                   |
|:--------------------------:|-----------------------------------------|-------------------------------------------------------|
|            `0`             | Internal error occurred                 | _Contact Backend                                      |
|            Dev_            |                                         |                                                       |
|            `1`             | Token invalid or not presented          | _Check if token is sent in `Authorization` header and |
| obtained with correct way_ |                                         |                                                       |
|            `2`             | Parse of json failed                    | _Check if websocket request is encoded to JSON        |
|         correctly_         |                                         |                                                       |
|            `3`             | Error occurred in sending message to MQ | _Check if Queue is correctly                          |
|        configured_         |                                         |                                                       |

## Deployment

### Branches

* Master: _Main version of Vertex_ 

