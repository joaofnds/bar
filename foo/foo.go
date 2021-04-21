package foo

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/joaofnds/bar/logger"
	"github.com/joaofnds/bar/tracing"
	"github.com/opentracing/opentracing-go"
)

var (
	FOO_SERVICE_URL = os.Getenv("FOO_SERVICE_URL")
)

func init() {
	if FOO_SERVICE_URL == "" {
		panic("I need FOO_SERVICE_URL to perform :cheems:")
	}

	logger.InfoLogger().Println("foo service initialized")
}

// CallFoo calls the foo service
func CallFoo(ctx context.Context) (string, error) {
	logger.DebugLogger().Println("CallFoo called")
	client := http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequest("GET", FOO_SERVICE_URL, nil)
	if err != nil {
		logger.ErrorLogger().Println("failed to build request, url: " + FOO_SERVICE_URL)
		return "", err
	}

	parentSpan := opentracing.SpanFromContext(ctx)
	span := opentracing.StartSpan("GET "+FOO_SERVICE_URL, opentracing.ChildOf(parentSpan.Context()))
	tracing.InjectRequestSpan(span, req)
	logger.DebugLogger().Println("will perform foo service call")
	resp, err := client.Do(req)
	logger.DebugLogger().Println("got foo service response")
	span.Finish()

	if err != nil {
		logger.ErrorLogger().Printf("failed to read body: %+v\n", err)
		return "", fmt.Errorf("failed to communicate with foo service: %w", err)
	}
	defer resp.Body.Close()

	logger.DebugLogger().Println("will read bofy")
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.ErrorLogger().Printf("failed to read body: %+v\n", err)
		return "", fmt.Errorf("failed to read body: %w", err)
	}

	logger.DebugLogger().Println("foo call went ok")

	return string(body), nil
}
