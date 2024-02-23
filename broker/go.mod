module github.com/DpodDani/broker

go 1.20

replace github.com/DpodDani/go-microservices-toolbox => ../toolbox

require (
	github.com/go-chi/chi/v5 v5.0.10
	github.com/go-chi/cors v1.2.1

)

require github.com/rabbitmq/amqp091-go v1.9.0

// generated this by running go get github.com/DpodDani/go-microservices-toolbox
// Note: must have added the "replace" statement from above first!
require github.com/DpodDani/go-microservices-toolbox v0.0.0-00010101000000-000000000000
