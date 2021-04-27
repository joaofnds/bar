package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/joaofnds/bar/config"
	"github.com/joaofnds/bar/foo_client"
	"github.com/joaofnds/bar/logger"
	"github.com/joaofnds/bar/tracing"
)

func main() {
	err := config.Parse()
	if err != nil {
		logger.ErrorLogger().Fatalf("failed to parse config: %v\n", err)
	}

	logger.InfoLogger().Println("Starting the application...")

	host, _ := os.Hostname()
	closer := tracing.InitTracer(host)
	defer closer.Close()

	fooService := foo_client.NewFooClient(config.FooServiceEndpoint(), 0)
	serviceName := config.ServiceName()

	http.HandleFunc("/", newFoohandler(serviceName, fooService))
	http.HandleFunc("/health", healthHandler)

	s := http.Server{Addr: ":80"}
	go func() {
		logger.ErrorLogger().Fatal(s.ListenAndServe())
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan

	logger.InfoLogger().Println("Shutdown signal received, exiting...")

	s.Shutdown(context.Background())
}
