package tracing

import (
	"fmt"
	"io"
	"net/http"

	"github.com/joaofnds/bar/config"
	opentracing "github.com/opentracing/opentracing-go"
	jaeger "github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-client-go/zipkin"
	"github.com/uber/jaeger-lib/metrics"
)

// initJaeger returns an instance of Jaeger Tracer that samples 100% of traces and logs all spans to stdout.
func InitTracer(serviceName string) io.Closer {
	collectorEndpoint := config.JaegerCollectorEndpoint()

	cfg := &jaegercfg.Configuration{
		ServiceName: serviceName,
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			CollectorEndpoint: collectorEndpoint,
		},
	}

	jMetricsFactory := metrics.NullFactory
	zipkinPropagator := zipkin.NewZipkinB3HTTPHeaderPropagator()

	closer, err := cfg.InitGlobalTracer(
		serviceName,
		jaegercfg.Logger(jaeger.StdLogger),
		jaegercfg.Metrics(jMetricsFactory),
		jaegercfg.Injector(opentracing.HTTPHeaders, zipkinPropagator),
		jaegercfg.Extractor(opentracing.HTTPHeaders, zipkinPropagator),
		jaegercfg.ZipkinSharedRPCSpan(true),
	)

	if err != nil {
		panic(fmt.Sprintf("ERROR: cannot init Jaeger: %v\n", err))
	}

	return closer
}

func spanCtxFromReq(tracer opentracing.Tracer, r *http.Request) (opentracing.SpanContext, error) {
	return tracer.Extract(
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(r.Header),
	)
}

func StartSpanFromReq(opName string, t opentracing.Tracer, r *http.Request) opentracing.Span {
	spanCtx, err := spanCtxFromReq(t, r)
	if err != nil {
		return opentracing.StartSpan(opName)
	}

	return opentracing.StartSpan(opName, opentracing.ChildOf(spanCtx))
}

func InjectRequestSpan(span opentracing.Span, request *http.Request) error {
	return span.Tracer().Inject(
		span.Context(),
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(request.Header))
}
