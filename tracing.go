package gin_opentelemetry_tracing

import (
	"context"
	"net/http"
	"runtime"
	"strconv"
	"strings"

	"go.opentelemetry.io/otel/api/kv/value"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/api/correlation"
	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/api/kv"
	"go.opentelemetry.io/otel/api/trace"
	"go.opentelemetry.io/otel/plugin/httptrace"
	"google.golang.org/grpc/codes"
)

const (
	TracingCtxKey = "tracing-context"
)

var (
	Caller  bool
	tracing = global.TraceProvider().Tracer("")
)

func SetComponentName(componentName string) {
	tracing = global.TraceProvider().Tracer(componentName)
}

func MiddlewareTracer() gin.HandlerFunc {
	return func(c *gin.Context) {
		attrs, entries, spanCtx := httptrace.Extract(c, c.Request)

		c.Request = c.Request.WithContext(correlation.ContextWithMap(c, correlation.NewMap(correlation.MapUpdate{
			MultiKV: entries,
		})))

		ctx, span := tracing.Start(
			trace.ContextWithRemoteSpanContext(c, spanCtx),
			c.HandlerName(),
			trace.WithAttributes(attrs...),
		)

		span.SetAttributes(kv.KeyValue{Key: "http.target", Value: value.String(c.Request.RequestURI)})
		span.SetAttributes(kv.KeyValue{Key: "http.host", Value: value.String(c.Request.URL.Host)})
		span.SetAttributes(kv.KeyValue{Key: "http.scheme", Value: value.String(c.Request.URL.Scheme)})
		span.SetAttributes(kv.KeyValue{Key: "http.flavor", Value: value.String(c.Request.Proto)})
		span.SetAttributes(kv.KeyValue{Key: "http.user_agent", Value: value.String(c.Request.UserAgent())})
		span.SetAttributes(httptrace.HTTPRemoteAddr.String(c.ClientIP()))

		defer span.End() // after all the other defers are completed.. finish the span

		c.Set(TracingCtxKey, ctx)
		c.Next()

		s := c.Writer.Status()
		span.SetAttributes(httptrace.HTTPStatus.Int(s))
		span.SetAttributes(kv.KeyValue{Key: "http.status_text", Value: value.String(http.StatusText(s))})
		span.SetStatus(mappingHTTPCodes(s), http.StatusText(s))
	}
}

func NewGinCtxSpan(c *gin.Context, name string) (ctx context.Context, span trace.Span) {
	var (
		spanName strings.Builder
	)

	if name != "" {
		spanName.WriteString(name)
	}

	if Caller {
		f, l := getFrame(1) // nolint:gomnd, gocritic
		spanName.WriteString(f)
		spanName.WriteString(":")
		spanName.WriteString(strconv.Itoa(l))
	} else {
		spanName.WriteString(c.HandlerName())
	}

	nCtx, b := c.Get(TracingCtxKey)

	if b {
		ctx = nCtx.(context.Context)

		return tracing.Start(ctx,
			spanName.String(),
		)
	}

	// If Trace not exist in ctx
	ctx = context.Background()

	opts := trace.WithAttributes(kv.KeyValue{Key: "http.target", Value: value.String(c.Request.RequestURI)},
		kv.KeyValue{Key: "http.host", Value: value.String(c.Request.URL.Host)},
		kv.KeyValue{Key: "http.scheme", Value: value.String(c.Request.URL.Scheme)},
		kv.KeyValue{Key: "http.flavor", Value: value.String(c.Request.Proto)},
		kv.KeyValue{Key: "http.user_agent", Value: value.String(c.Request.UserAgent())},
		httptrace.HTTPRemoteAddr.String(c.ClientIP()))

	return tracing.Start(ctx,
		spanName.String(),
		opts,
	)
}

// nolint: gomnd, gocyclo, gocritic
/*
 Reference https://github.com/open-telemetry/opentelemetry-specification/blob/master/specification/data-http.md
 and
 https://github.com/census-instrumentation/opencensus-specs/blob/master/trace/HTTP.md#mapping-from-http-status-codes-to-trace-status-codes
*/
func mappingHTTPCodes(c int) (r codes.Code) {
	switch {
	case c < 100:
		r = codes.Unknown
	case c <= 399:
		r = codes.OK
	case c == 400:
		r = codes.InvalidArgument
	case c == 401:
		r = codes.Unauthenticated
	case c == 403:
		r = codes.PermissionDenied
	case c == 404:
		r = codes.NotFound
	case c == 429:
		r = codes.ResourceExhausted
	case c < 500:
		r = codes.InvalidArgument
	case c == 501:
		r = codes.Unimplemented
	case c == 503:
		r = codes.Unavailable
	case c == 504:
		r = codes.DeadlineExceeded
	case c < 600:
		r = codes.Internal
	}

	return
}

// Get runtime.Frame
// nolint: gomnd, gocritic
func getFrame(skipFrames int) (string, int) {
	// We need the frame at index skipFrames+2, since we never want runtime.Callers and getFrame
	targetFrameIndex := skipFrames + 2

	// Set size to targetFrameIndex+2 to ensure we have room for one more caller than we need
	programCounters := make([]uintptr, targetFrameIndex+2)
	n := runtime.Callers(0, programCounters)

	frame := runtime.Frame{Function: "unknown"}

	if n > 0 {
		frames := runtime.CallersFrames(programCounters[:n])

		for more, frameIndex := true, 0; more && frameIndex <= targetFrameIndex; frameIndex++ {
			var frameCandidate runtime.Frame
			frameCandidate, more = frames.Next()

			if frameIndex == targetFrameIndex {
				frame = frameCandidate
			}
		}
	}

	return frame.Function, frame.Line
}
