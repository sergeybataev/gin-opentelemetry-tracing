package main

import (
	"fmt"
	"net/http"
	"path"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	tracing "github.com/sergeybataev/gin-opentelemetry-tracing"
	"github.com/sergeybataev/gin-opentelemetry-tracing/example"

	"github.com/jinzhu/configor"
)

func main() {
	err := configor.New(&configor.Config{Debug: true, Verbose: true}).
		Load(&example.Config, path.Join("./example/config.yml"))

	log.Debug().Msgf("%v", err)
	log.Debug().Msgf("%v", example.Config)

	// Store trace into Jaeger
	/*	fn := example.NewJaegerTracing(example.Config.Jaeger.Endpoint, "Gateway")
		defer fn()*/
	// Print trace info into Stdout
	example.NewStdoutTracing("Gateway")

	// Set print function:line
	tracing.Caller = true

	// Set up routes
	r := gin.Default()
	r.Use(tracing.MiddlewareTracer())

	// HealthCheck
	// Ping test
	r.GET("/healthcheck", HealthCheck)

	_ = r.Run(fmt.Sprintf("%s:%d", example.Config.Gateway.Host, example.Config.Gateway.Port))
}

func HealthCheck(c *gin.Context) {
	_, span := tracing.NewGinCtxSpan(c, "healthcheck")
	defer span.End()

	c.String(http.StatusOK, "ok")
}
