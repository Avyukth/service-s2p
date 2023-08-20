package web

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"syscall"
	"time"

	"github.com/google/uuid"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"

	"github.com/dimfeld/httptreemux/v5"
)

type Handler func(ctx context.Context, w http.ResponseWriter, r *http.Request) error

type App struct {
	mux *httptreemux.ContextMux
	otmux http.Handler
	shutdown chan os.Signal
	mw       []Middleware
}

func NewApp(shutdown chan os.Signal, mw ...Middleware) *App {
	
	mux := httptreemux.NewContextMux()
	return &App{
		mux: mux,
		otmux: otelhttp.NewHandler(mux, "request"),
		shutdown:   shutdown,
		mw:         mw,
	}
}

func (a *App) SignalShutdown() {
	a.shutdown <- syscall.SIGTERM
}

func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request){
	a.otmux.ServeHTTP(w,r)
}

func (a *App) Handle(method string, group string, path string, handler Handler, mw ...Middleware) {

	handler = wrapMiddleware(mw, handler)

	handler = wrapMiddleware(a.mw, handler)


	h := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

	span:= trace.SpanFromContext(ctx)
		v := Values{
			TraceID: span.SpanContext().SpanID().String(),
			Now:     time.Now(),
		}

		ctx = context.WithValue(ctx, key, &v)

		if err := handler(ctx, w, r); err != nil {
			a.SignalShutdown()
			return
		}

	}

	finalPath := path
	if group != "" {
		finalPath = "/" + group + path
	}
	a.mux.Handle(method, finalPath, h)
}


// =============================================================================

// startTracing configure open telemetry to be used with Grafana Tempo.

func startTracing(serviceName string, reporterURI string, probability float64) (*trace.TracerProvider, error){

	exporter, err := otlptrace.New(
		context.Background(),
		otlptracegrpc.NewClient(
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint(reporterURI),
		),
	)

	if err != nil {
		return nil, fmt.Errorf("creating new otel exporter: %w", err)
	}

	traceProvider := trace.NewNoopTracerProvider(
		trace.WithSampler(trace.TraceIDRatioBased(probability),
	trace.WithBatcher(exporter, trace.WithMaxExportBatchSize(trace.DefaultMaxExportBatchSize),
	trace.WithBatchTimeout(trace.DefaultScheduleDelay*time.Millisecond),
	trace.WithMaxExportBatchSize(trace.DefaultMaxExportBatchSize),),
trace.WithResource(resource.NewWithAttributes(semconv.SchemaURL, semconv.ServiceNameKey.String(serviceName),),),
	
	

}
