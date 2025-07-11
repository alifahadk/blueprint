module github.com/blueprint-uservices/cldrel_course/luggageshare/wiring

go 1.22.1

require (
	github.com/blueprint-uservices/blueprint/blueprint v0.0.0-20250511234457-4e6e55620828
	github.com/blueprint-uservices/blueprint/plugins v0.0.0-20250511234457-4e6e55620828
	github.com/blueprint-uservices/cldrel_course/luggageshare/e2etests v0.0.0
	github.com/blueprint-uservices/cldrel_course/luggageshare/workflow v0.0.0
	github.com/blueprint-uservices/cldrel_course/luggageshare/workload v0.0.0
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/blueprint-uservices/blueprint/runtime v0.0.0-20250324124051-1e5478a0a8a6 // indirect
	github.com/go-logr/logr v1.4.1 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-sql-driver/mysql v1.8.1 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/jmoiron/sqlx v1.4.0 // indirect
	github.com/klauspost/compress v1.17.8 // indirect
	github.com/montanaflynn/stats v0.7.1 // indirect
	github.com/otiai10/copy v1.14.0 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.1.2 // indirect
	github.com/xdg-go/stringprep v1.0.4 // indirect
	github.com/youmark/pkcs8 v0.0.0-20240424034433-3c2c7870ae76 // indirect
	go.mongodb.org/mongo-driver v1.15.0 // indirect
	go.opentelemetry.io/otel v1.26.0 // indirect
	go.opentelemetry.io/otel/exporters/jaeger v1.17.0 // indirect
	go.opentelemetry.io/otel/exporters/stdout/stdoutmetric v1.26.0 // indirect
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.26.0 // indirect
	go.opentelemetry.io/otel/metric v1.26.0 // indirect
	go.opentelemetry.io/otel/sdk v1.26.0 // indirect
	go.opentelemetry.io/otel/sdk/metric v1.26.0 // indirect
	go.opentelemetry.io/otel/trace v1.26.0 // indirect
	golang.org/x/crypto v0.22.0 // indirect
	golang.org/x/exp v0.0.0-20240416160154-fe59bbe5cc7f // indirect
	golang.org/x/mod v0.17.0 // indirect
	golang.org/x/sync v0.7.0 // indirect
	golang.org/x/sys v0.19.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	golang.org/x/tools v0.20.0 // indirect
	gonum.org/v1/gonum v0.15.1 // indirect
)

replace github.com/blueprint-uservices/cldrel_course/luggageshare/workflow => ../workflow

replace github.com/blueprint-uservices/cldrel_course/luggageshare/workload => ../workload

replace github.com/blueprint-uservices/cldrel_course/luggageshare/e2etests => ../e2etests
