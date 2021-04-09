package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/joaofnds/bar/foo"
)

func main() {
	log.Println("Starting the application...")

	var host string

	if hostname, err := os.Hostname(); err == nil {
		host = hostname
	} else {
		host = "some server"
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if response, err := foo.CallFoo(); err != nil {
			log.Println(err)
			fmt.Fprintln(w, "Hello from "+host+", I failed to contact foo service")
		} else {
			fmt.Fprintln(w, "Hello from "+host+", here's what foo service said: "+response)
		}
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	})

	s := http.Server{Addr: ":80"}
	go func() {
		log.Fatal(s.ListenAndServe())
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan

	log.Println("Shutdown signal received, exiting...")

	s.Shutdown(context.Background())
}
