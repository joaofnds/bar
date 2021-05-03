package main

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_healthHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "", nil)
	if err != nil {
		t.Fatalf("could no create request: %v", err)
	}
	rec := httptest.NewRecorder()

	healthHandler(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status OK, got: %v", res.StatusCode)
	}

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("failed to read body")
	}

	if len(b) != 0 {
		t.Errorf("body must be empty, got: %s", b)
	}
}

func Test_newFooHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "", nil)
	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}

	tt := []struct {
		name        string
		fooResponse string
	}{
		{"foo", "hi from foo"},
		{"bar", "hi from bar"},
		{"qux", "I'm just qux"},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()

			fooSvcMock := FooServiceMock{tc.fooResponse, nil}
			newFoohandler(tc.name, fooSvcMock)(rec, req)

			res := rec.Result()
			defer res.Body.Close()

			if res.StatusCode != http.StatusOK {
				t.Errorf("expected status OK, got: %v", res.Status)
			}

			b, err := ioutil.ReadAll(res.Body)
			if err != nil {
				t.Fatalf("failed to read body: %v", err)
			}

			expectedResponse := fmt.Sprintf("Hello from %s, here's what foo service said: %s", tc.name, tc.fooResponse)

			if string(b) != expectedResponse {
				t.Fatalf("expected response to be %q, but got: %q", expectedResponse, b)
			}
		})
	}
}

func Test_newFooHandlerWhenServiceFails(t *testing.T) {
	req, err := http.NewRequest("GET", "", nil)
	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}

	tt := []struct {
		name   string
		fooErr error
	}{
		{"foo", errors.New("oops")},
		{"bar", errors.New("my bad")},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()

			fooSvcMock := FooServiceMock{"", tc.fooErr}
			newFoohandler(tc.name, fooSvcMock)(rec, req)

			res := rec.Result()
			defer res.Body.Close()

			if res.StatusCode != http.StatusFailedDependency {
				t.Errorf("expected status FailedDependency, got: %v", res.Status)
			}

			b, err := ioutil.ReadAll(res.Body)
			if err != nil {
				t.Fatalf("failed to read body: %v", err)
			}

			expectedResponse := fmt.Sprintf("Hello from super %s, I failed to contact foo service", tc.name)

			if string(b) != expectedResponse {
				t.Fatalf("expected response to be %q, but got: %q", expectedResponse, b)
			}
		})
	}
}

type FooServiceMock struct {
	str string
	err error
}

func (s FooServiceMock) CallFoo(_ context.Context) (string, error) {
	return s.str, s.err
}
