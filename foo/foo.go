package foo

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

var (
	FOO_SERVICE_URL = os.Getenv("FOO_SERVICE_URL")
)

func init() {
	if FOO_SERVICE_URL == "" {
		panic("I need FOO_SERVICE_URL to perform :cheems:")
	}
}

// CallFoo calls the foo service
func CallFoo() (string, error) {
	client := http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(FOO_SERVICE_URL)
	if err != nil {
		return "", fmt.Errorf("failed to communicate with foo service: %w", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read body: %w", err)
	}

	return string(body), nil
}
