package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joaofnds/bar/foo"
	"github.com/joaofnds/bar/logger"
)

func main() {
	logger.InfoLogger().Println("Starting the application...")

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
	defer trackTime(time.Now(), "rootHandler")
	logger.InfoLogger().Printf("started rootHandler: %v\n", r.URL)

	host, _ := os.Hostname()

	response, err := foo.CallFoo()
	if err != nil {
		logger.ErrorLogger().Printf("failed to call foo service: %+v\n", err)
		fmt.Fprintln(w, "Hello from "+host+", I failed to contact foo service")

		return
	}

	fmt.Fprintln(w, "Hello from "+host+", here's what foo service said: "+response)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	defer trackTime(time.Now(), "healthHandler")
	w.WriteHeader(200)
}

func trackTime(start time.Time, funcName string) {
	elapsed := time.Since(start)
	logger.InfoLogger().Printf("finished %s in %s\n", funcName, elapsed)
}
