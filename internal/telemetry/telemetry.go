package telemetry

import (
	"context"
	"errors"
	"io"
	"os"
	"time"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewResource(serviceName, serviceVersion string) *resource.Resource {
	return resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName(serviceName),
		semconv.ServiceNamespace("ssemu"),
		semconv.ServiceVersion(serviceVersion),
	)
}

func NewTraceProvider(ctx context.Context, res *resource.Resource) (*trace.TracerProvider, error) {
	var err error
	var exp trace.SpanExporter
	collectorEndpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if len(collectorEndpoint) > 0 {
		exp, err = newGrpcTraceExporter(ctx, collectorEndpoint)
	} else {
		exp, err = stdouttrace.New(
			stdouttrace.WithWriter(io.Discard),
		)
	}
	if err != nil {
		return nil, err
	}
	return trace.NewTracerProvider(
		trace.WithResource(res),
		trace.WithBatcher(exp, trace.WithBatchTimeout(5*time.Second)),
	), nil
}

func NewMeterProvider(ctx context.Context, res *resource.Resource) (*metric.MeterProvider, error) {
	exp, err := prometheus.New(
		prometheus.WithNamespace("ssemu"),
	)
	if err != nil {
		return nil, err
	}
	return metric.NewMeterProvider(
		metric.WithResource(res),
		metric.WithReader(exp),
	), nil
}

func newGrpcTraceExporter(ctx context.Context, collectorEndpoint string) (trace.SpanExporter, error) {
	conn, err := openGrpcConn(ctx, collectorEndpoint)
	if err != nil {
		return nil, err
	}
	return otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
}

func openGrpcConn(ctx context.Context, collectorEndpoint string) (*grpc.ClientConn, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	conn, err := grpc.DialContext(ctx, collectorEndpoint,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if errors.Is(err, context.DeadlineExceeded) {
		return nil, errors.New("telemetry grpc dial timeout")
	}
	return conn, nil
}
