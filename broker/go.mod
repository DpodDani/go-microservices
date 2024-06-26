module github.com/DpodDani/broker

go 1.20

replace github.com/DpodDani/go-microservices-toolbox => ../toolbox

require (
	github.com/go-chi/chi/v5 v5.0.10
	github.com/go-chi/cors v1.2.1

)

require github.com/rabbitmq/amqp091-go v1.9.0

require (
	github.com/golang/protobuf v1.5.3 // indirect
	golang.org/x/net v0.20.0 // indirect
	golang.org/x/sys v0.16.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240123012728-ef4313101c80 // indirect
	google.golang.org/grpc v1.62.0 // indirect
	google.golang.org/protobuf v1.32.0 // indirect
)

// generated this by running go get github.com/DpodDani/go-microservices-toolbox
// Note: must have added the "replace" statement from above first!
require github.com/DpodDani/go-microservices-toolbox v0.0.0-00010101000000-000000000000
