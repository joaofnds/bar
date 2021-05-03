package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/joaofnds/bar/logger"
	"github.com/joaofnds/bar/tracing"
	"github.com/opentracing/opentracing-go"
)

type fooCaller interface {
	CallFoo(context.Context) (string, error)
}

func newFoohandler(name string, fooCaller fooCaller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		span := tracing.StartSpanFromReq("rootHandler", opentracing.GlobalTracer(), r)
		defer span.Finish()
		ctx := opentracing.ContextWithSpan(context.Background(), span)

		response, err := fooCaller.CallFoo(ctx)
		if err != nil {
			logger.ErrorLogger().Printf("failed to call foo service: %+v\n", err)
			w.WriteHeader(http.StatusFailedDependency)
			fmt.Fprint(w, "Hello from ["+name+"], I failed to contact foo service")

			return
		}

		fmt.Fprint(w, "Hello from "+name+", here's what foo service said: "+response)
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	span := tracing.StartSpanFromReq("healthHandler", opentracing.GlobalTracer(), r)
	defer span.Finish()

	w.WriteHeader(200)
}
