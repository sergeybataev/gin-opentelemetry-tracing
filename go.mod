module github.com/sergeybataev/gin-opentelemetry-tracing

go 1.14

require (
	github.com/gin-gonic/gin v1.6.3
	github.com/jinzhu/configor v1.2.0
	github.com/rs/zerolog v1.18.0
	go.opentelemetry.io/otel v0.5.0
	go.opentelemetry.io/otel/exporters/otlp v0.5.0
	go.opentelemetry.io/otel/exporters/trace/jaeger v0.5.0
	google.golang.org/grpc v1.29.1
)
