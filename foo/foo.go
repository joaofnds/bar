package foo

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/joaofnds/bar/logger"
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
func CallFoo() (string, error) {
	client := http.Client{Timeout: 5 * time.Second}

	start := time.Now()

	resp, err := client.Get(FOO_SERVICE_URL)

	elapsed := time.Since(start)
	logger.InfoLogger().Printf("finished foo service call in %s", elapsed)

	if err != nil {
		logger.ErrorLogger().Printf("failed to read body: %+v\n", err)
		return "", fmt.Errorf("failed to communicate with foo service: %w", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.ErrorLogger().Printf("failed to read body: %+v\n", err)
		return "", fmt.Errorf("failed to read body: %w", err)
	}

	return string(body), nil
}
