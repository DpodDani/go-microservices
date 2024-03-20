FROM golang:1.20-alpine as builder

RUN mkdir /app

COPY ./front-end /app

WORKDIR /app

RUN CGO_ENABLED=0 go build -o frontApp ./cmd/web

RUN chmod + /app/frontApp

FROM alpine:latest

RUN mkdir /app

COPY --from=builder /app/frontApp /app

CMD ["/app/frontApp"]