FRONT_END_BINARY=frontApp
BROKER_BINARY=brokerApp
LOGGER_BINARY=loggerApp
AUTH_BINARY=authApp
MAIL_BINARY=mailApp
LISTENER_BINARY=listenerApp

# TODO: Run SQL scripts in Postgres docker container to init DB!

## up: starts all containers in the background without forcing build
docker-up:
	@echo "Starting Docker images..."
	docker-compose up -d
	@echo "Docker images started!"

## build_up: stops docker-compose (if running), builds all projects and starts docker compose
docker-build-up: build-broker build-auth build-logger build-mail build-listener
	@echo "Stopping docker images (if running...)"
	docker-compose down
	@echo "Building (when required) and starting docker images..."
	docker-compose up --build -d
	@echo "Docker images built and started!"

## rebuild: stops docker-compose (if running) and removes all images,
## builds all projects and starts docker compose
docker-rebuild: build-broker build-auth build-logger build-mail build-listener
	@echo "Stopping docker images (if running...)"
	docker-compose down --rmi all
	@echo "Building (when required) and starting docker images..."
	docker-compose up --build -d
	@echo "Docker images built and started!"

## down: stop docker compose
docker-down:
	@echo "Stopping docker compose..."
	docker-compose down
	@echo "Done!"

## build_broker: builds the broker binary as a linux executable
build-broker:
	@echo "Building broker binary..."
	cd broker && env GOOS=linux CGO_ENABLED=0 go build -o ${BROKER_BINARY} ./cmd/api
	@echo "Done!"

## build_auth: builds the auth binary as a linux executable
build-auth:
	@echo "Building auth binary..."
	cd auth && env GOOS=linux CGO_ENABLED=0 go build -o ${AUTH_BINARY} ./cmd/api
	@echo "Done!"

## build_logger: builds the logger binary as a linux executable
build-logger:
	@echo "Building logger binary..."
	cd logger && env GOOS=linux CGO_ENABLED=0 go build -o ${LOGGER_BINARY} ./cmd/api
	@echo "Done!"

# build_mail: builds the mail binary as a linux executable
build-mail:
	@echo "Building mail binary..."
	cd mail && env GOOS=linux CGO_ENABLED=0 go build -o ${MAIL_BINARY} ./cmd/api
	@echo "Done!"

# build_listener: builds the listener binary as a linux executable
build-listener:
	@echo "Building listener binary..."
	cd listener && env GOOS=linux CGO_ENABLED=0 go build -o ${LISTENER_BINARY} ./
	@echo "Done!"

## build_front: builds the frone end binary
build-frontend:
	@echo "Building front end binary..."
	cd front-end && env CGO_ENABLED=0 go build -o ${FRONT_END_BINARY} ./cmd/web
	@echo "Done!"

## start: starts the front end
start-frontend: build-frontend
	@echo "Starting front end"
	cd front-end && ./${FRONT_END_BINARY} &

## stop: stop the front end
stop-frontend:
	@echo "Stopping front end..."
	@-pkill -SIGTERM -f "./${FRONT_END_BINARY}"
	@echo "Stopped front end!"


## Compile proto files

root_dirs := logger broker
compile_cmd := protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative *.proto

compile-proto-files:
	@echo "Compiling proto files for services: $(root_dirs)"
	# For each "svc" in services, run the protoc compile command
	$(foreach svc,$(root_dirs),$$(cd $(svc)/proto && $$($(compile_cmd))))
	@echo "Finished compiling proto files"


## Upload images to Docker hub
# Note: broker service is handled differently
services := logger auth mail listener

upload-images:
	@echo "Uploading images to Docker Hub"
	# build broker image
	docker build -f broker/broker.dockerfile -t dnamufetha/broker-service:1.0.0 .
	# build images for other services
	$(foreach svc,$(services), $$(cd $(svc) && $$(docker build -f $(svc).dockerfile -t dnamufetha/$(svc)-service:1.0.0 .)))

	# log into Docker Hub
	docker login -u dnamufetha --password-stdin < .docker_pat.token

	# upload images to Docker Hub
	docker push dnamufetha/broker-service:1.0.0
	docker push dnamufetha/logger-service:1.0.0
	docker push dnamufetha/auth-service:1.0.0
	docker push dnamufetha/mail-service:1.0.0
	docker push dnamufetha/listener-service:1.0.0