module blueprint/testclients

go 1.22.1

toolchain go1.24.2

replace github.com/blueprint-uservices/blueprint/runtime => ../runtime

replace github.com/blueprint-uservices/cldrel_course/luggageshare/e2etests => ../e2etests

replace github.com/blueprint-uservices/cldrel_course/luggageshare/workflow => ../workflow

require (
	github.com/blueprint-uservices/blueprint/runtime v0.0.0-20250324124051-1e5478a0a8a6
	github.com/blueprint-uservices/cldrel_course/luggageshare/workflow v0.0.0-00010101000000-000000000000
)

require (
	github.com/google/uuid v1.6.0 // indirect
	go.mongodb.org/mongo-driver v1.15.0 // indirect
	go.opentelemetry.io/otel v1.26.0 // indirect
	go.opentelemetry.io/otel/metric v1.26.0 // indirect
	go.opentelemetry.io/otel/trace v1.26.0 // indirect
	golang.org/x/exp v0.0.0-20240416160154-fe59bbe5cc7f // indirect
)
