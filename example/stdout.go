package example

import (
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel/api/global"
	tracestdout "go.opentelemetry.io/otel/exporters/trace/stdout"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"

	tracing "github.com/sergeybataev/gin-opentelemetry-tracing"
)

// initTracer creates a new trace provider instance and registers it as global trace provider.
func NewStdoutTracing(componentName string) {
	exp, err := tracestdout.NewExporter(tracestdout.Options{PrettyPrint: false})
	if err != nil {
		log.Fatal().Err(err).Msgf("failed to initialize trace stdout exporter")
		return
	}

	// For demoing purposes, always sample. In a production application, you should
	// configure this to a trace.ProbabilitySampler set at the desired
	// probability.
	tp, err := sdktrace.NewProvider(
		sdktrace.WithConfig(sdktrace.Config{DefaultSampler: sdktrace.AlwaysSample()}),
		sdktrace.WithSyncer(exp))
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	global.SetTraceProvider(tp)

	tracing.SetComponentName(componentName)
}
