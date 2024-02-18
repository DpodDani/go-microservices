# Micro-services in Golang

## Connecting to resources

### App UI

Navigate to http://localhost:80

### MongoDB

#### Connecting to MongoDB

```shell
mongo "mongodb://admin:password@localhost:27017/logs?authSource=admin"
```

#### Querying MongoDB

```shell
db.logs.find({})
```

### RabbitMQ

Navigate to http://localhost:15672/

### MailHog

Navigate to http://localhost:8025/

## Development

### Building service

```shell
make build-{SERVICE_NAME}
```

### Create docker image and start container for single service

```shell
docker compose up --build -d {SERVICE_NAME}-service
```

### Create docker images and start containers for ALL services

```shell
make docker-build-up
```

### Start/stop containers

```shell
make docker-up # run containers for all services
make docker-down # stop containers for all services (without deleting images)
```

### Compile proto files

1. Navigate to directory containing `proto` files, then run:

```shell
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative <PROTO_FILE>
```