package foo_client_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/joaofnds/bar/foo_client"
)

func TestCallFoo(t *testing.T) {
	server := newTestServer()
	defer server.Close()

	tt := []struct {
		name string
		path string
		err  error
		resp string
	}{
		{name: "success", path: "/success", err: nil, resp: "foo"},
		{name: "failure", path: "/failure", err: foo_client.ServerError, resp: ""},
		{name: "timeout", path: "/timeout", err: foo_client.DeadlineExceeded, resp: ""},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			svc := foo_client.NewFooClient(server.URL+tc.path, 10*time.Millisecond)
			resp, err := svc.CallFoo(context.Background())

			if resp != tc.resp {
				t.Errorf("expected response %#v, got: %#v", tc.resp, resp)
			}

			if err != tc.err {
				t.Errorf("expected err to contain %#v, but got: %#v", tc.err, err)
			}
		})
	}
}

func newTestServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/success":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("foo"))
		case "/failure":
			w.WriteHeader(http.StatusInternalServerError)
		case "/timeout":
			time.Sleep(time.Second)
			w.WriteHeader(http.StatusOK)
		}
	}))
}
