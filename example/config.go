package example

var Config = struct {
	Gateway struct {
		Host string
		Port uint `default:"8080"`
	}

	Service1 struct {
		Host string
		Port uint `default:"8081"`
	}

	Service2 struct {
		Host string
		Port uint `default:"8082"`
	}

	Service3 struct {
		Host string
		Port uint `default:"8083"`
	}

	Jaeger struct {
		Endpoint string
	}
	Opentelemetry struct {
		Endpoint string
	}
}{}
