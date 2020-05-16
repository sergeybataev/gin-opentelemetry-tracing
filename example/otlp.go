package example

import (
	"time"

	"github.com/rs/zerolog/log"
	tracing "github.com/sergeybataev/gin-opentelemetry-tracing"
	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/exporters/otlp"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func NewOptlTracing(endpoint, componentName string) func() {

	exporter, err := otlp.NewExporter(
		otlp.WithInsecure(),
		otlp.WithAddress(endpoint),
		otlp.WithReconnectionPeriod(50*time.Millisecond))
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
		exporter.Stop()
	}
}
