package main

import (
	"net/http/httptest"
	"path"
	"testing"

	"go.opentelemetry.io/otel/api/kv/value"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/configor"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel/api/kv"
	"go.opentelemetry.io/otel/api/trace"

	tracing "github.com/sergeybataev/gin-opentelemetry-tracing"
	"github.com/sergeybataev/gin-opentelemetry-tracing/example"
)

/*
func TestHealthCheck(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	HealthCheck(c)

	assert.Equal(t, 200, w.Result().Status) // or what value you need it to be
}*/

func BenchmarkHealthCheck(b *testing.B) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request = httptest.NewRequest("GET", "http://test.com", nil)

	c.Request.RequestURI = "test"
	c.Request.URL.Host = "testHost"
	c.Request.URL.Scheme = "testScheme"
	c.Request.Proto = "testProto"

	_, span := tracing.NewGinCtxSpan(c, "HandlerName Pong")
	span.End()

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		HealthCheck(c)
	}

}

func BenchmarkAttributes(b *testing.B) {
	var opts trace.StartOption

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {

		opts = trace.WithAttributes(kv.KeyValue{
			Key:   "http.target",
			Value: value.String("http_uri"),
		})
	}
	_ = opts
}

func init() {
	err := configor.Load(&example.Config, path.Join("./example/config.yml"))

	log.Debug().Msgf("%v", err)
	log.Debug().Msgf("%v", example.Config)

	// Set print function:line
	tracing.Caller = false

}
