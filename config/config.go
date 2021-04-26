package config

import (
	"os"
)

func Parse() error {
	return nil
}

func JaegerCollectorEndpoint() string {
	return os.Getenv("JAEGER_COLLECTOR_ENDPOINT")
}

func FooServiceEndpoint() string {
	return os.Getenv("FOO_SERVICE_URL")
}
