package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/joaofnds/bar/config"
	"github.com/joaofnds/bar/foo"
	"github.com/joaofnds/bar/logger"
	"github.com/joaofnds/bar/tracing"
	"github.com/opentracing/opentracing-go"
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

	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/health", healthHandler)

	s := http.Server{Addr: ":3000"}
	go func() {
		logger.ErrorLogger().Fatal(s.ListenAndServe())
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan

	logger.InfoLogger().Println("Shutdown signal received, exiting...")

	s.Shutdown(context.Background())
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	span := tracing.StartSpanFromReq("rootHandler", opentracing.GlobalTracer(), r)
	defer span.Finish()

	ctx := opentracing.ContextWithSpan(context.Background(), span)

	host, _ := os.Hostname()

	response, err := foo.CallFoo(ctx)
	if err != nil {
		logger.ErrorLogger().Printf("failed to call foo service: %+v\n", err)
		fmt.Fprintln(w, "Hello from "+host+", I failed to contact foo service")

		return
	}

	fmt.Fprintln(w, "Hello from "+host+", here's what foo service said: "+response)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	span := tracing.StartSpanFromReq("healthHandler", opentracing.GlobalTracer(), r)
	defer span.Finish()

	w.WriteHeader(200)
}
