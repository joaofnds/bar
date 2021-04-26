package foo_client

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/joaofnds/bar/logger"
	"github.com/joaofnds/bar/tracing"
	"github.com/opentracing/opentracing-go"
)

var (
	ServerError      = errors.New("server exploded")
	DeadlineExceeded = errors.New("deadline exceeded")
)

type FooClient struct {
	endpoint string
	timeout  time.Duration
}

func NewFooClient(endpoint string, timeout time.Duration) *FooClient {
	return &FooClient{endpoint: endpoint, timeout: timeout}
}

// CallFoo calls the foo service
func (s *FooClient) CallFoo(ctx context.Context) (string, error) {
	client := http.Client{Timeout: s.timeout}
	req, err := http.NewRequest("GET", s.endpoint, nil)
	if err != nil {
		logger.ErrorLogger().Println("failed to build request, url: " + s.endpoint)
		return "", err
	}

	parentSpan := opentracing.SpanFromContext(ctx)
	if parentSpan != nil {
		span := opentracing.StartSpan("GET "+s.endpoint, opentracing.ChildOf(parentSpan.Context()))
		tracing.InjectRequestSpan(span, req)
		defer span.Finish()
	}
	resp, err := client.Do(req)

	if err != nil {
		switch err := err.(type) {
		case *url.Error:
			if err.Timeout() {
				logger.ErrorLogger().Printf("failed to perform request: %+v\n", err)
				return "", DeadlineExceeded
			} else {
				return "", fmt.Errorf("unknown error: %w", err)
			}
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", ServerError
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.ErrorLogger().Printf("failed to read body: %+v\n", err)
		return "", fmt.Errorf("failed to read body: %w", err)
	}

	return string(body), nil
}
