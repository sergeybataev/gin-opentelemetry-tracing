package example

import (
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/api/kv"
	"go.opentelemetry.io/otel/exporters/trace/jaeger"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"

	tracing "github.com/sergeybataev/gin-opentelemetry-tracing"
)

func NewJaegerTracing(endpoint, componentName string) func() {
	exporter, err := jaeger.NewRawExporter(
		jaeger.WithCollectorEndpoint(endpoint),
		jaeger.WithProcess(jaeger.Process{
			ServiceName: componentName,
			Tags: []kv.KeyValue{
				kv.String("exporter", "jaeger"),
				kv.Float64("float", 312.23),
			},
		}),
		//		jaeger.
	)
	if err != nil {
		log.Fatal().Err(err)
	}

	tp, err := sdktrace.NewProvider(
		sdktrace.WithConfig(sdktrace.Config{DefaultSampler: sdktrace.AlwaysSample()}),
		sdktrace.WithSyncer(exporter))
	if err != nil {
		log.Fatal().Err(err)
	}

	global.SetTraceProvider(tp)

	tracing.SetComponentName(componentName)

	return func() {
		exporter.Flush()
	}
}
