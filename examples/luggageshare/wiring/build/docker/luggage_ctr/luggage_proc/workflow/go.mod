module github.com/blueprint-uservices/cldrel_course/luggageshare/workflow

go 1.22.1

require (
	github.com/blueprint-uservices/blueprint/runtime v0.0.0-20250324124051-1e5478a0a8a6
	github.com/google/uuid v1.6.0
)

require (
	go.mongodb.org/mongo-driver v1.15.0 // indirect
	go.opentelemetry.io/otel v1.26.0 // indirect
	go.opentelemetry.io/otel/metric v1.26.0 // indirect
	go.opentelemetry.io/otel/trace v1.26.0 // indirect
	golang.org/x/exp v0.0.0-20240416160154-fe59bbe5cc7f // indirect
)

replace github.com/blueprint-uservices/blueprint/runtime => ../runtime

replace blueprint/goproc/luggage_proc => ../luggage_proc
