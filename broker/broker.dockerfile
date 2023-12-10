# base golang image
FROM golang:1.20-alpine as builder

RUN mkdir /app /toolbox

COPY ../toolbox /toolbox

COPY ./broker /app

WORKDIR /app

# not using any C libraries, so don't need CGO enabled
RUN CGO_ENABLED=0 go build -o brokerApp ./cmd/api

RUN chmod +x /app/brokerApp

# tiny image
FROM alpine:latest

RUN mkdir /app

COPY --from=builder /app/brokerApp /app

CMD ["/app/brokerApp"]