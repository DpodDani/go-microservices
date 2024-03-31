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

### Add backend host

1. Edit `/etc/hosts` file (requires `sudo`)

```bash
##
# Host Database
#
# localhost is used to configure the loopback interface
# when the system is booting.  Do not change this entry.
##
127.0.0.1       localhost backend # added the word 'backend' here
255.255.255.255 broadcasthost
::1             localhost backend # added the word 'backend' here
```

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

> Alternatively, run `make compile-proto-files` in project root

## Docker swarm

### For new projects, you have to initialise Docker swarm like so:

```shell
docker swarm init
```

### View token to have manager/worker node join swarm:

```shell
docker swarm join-token [manager|worker]
```

### Scaling services up/down

```shell
docker service scale <SERVICE_NAME>=<REPLICA_NUM>
```

### Updating service in docker swarm

1. Create new tagged version of service
2. Push newer versioned image to Docker Hub
3. Scale the service to at least 2 instances (to prevent downtime when updating image)
4. Update image in swarm to use newer version:
   * Note: can also use to rollback to stable version (when bug found in new version)
```shell
docker service update --image <NEW_IMAGE> <SERVICE_NAME>
```

### Remove services from swarm

```shell
docker stack rm <SWARM_NAME>
```

### Make node leave swarm

```shell
docker swarm leave
```