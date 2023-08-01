package o11y

import (
	"context"
	"log"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"google.golang.org/grpc/credentials"

	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

type Opts struct {
	ServiceName     string
	CollectorHeader map[string]string
	CollectorUrl    string
	InsecureMode    bool

	// Bundler specific attributes
	ChainID *big.Int
	Address common.Address
}

func initResources(opts *Opts) *resource.Resource {
	resources, err := resource.New(
		context.Background(),
		resource.WithAttributes(
			attribute.String("service.name", opts.ServiceName),
			attribute.String("library.language", "go"),
			attribute.String("bundler.address", opts.Address.Hex()),
			attribute.Int64("bundler.chain_id", opts.ChainID.Int64()),
		),
	)
	if err != nil {
		log.Fatal(err)
	}

	return resources
}

func IsEnabled(serviceName string) bool {
	return len(serviceName) > 0
}

func InitTracer(opts *Opts) func() {
	secureOption := otlptracegrpc.WithTLSCredentials(credentials.NewClientTLSFromCert(nil, ""))
	if opts.InsecureMode {
		secureOption = otlptracegrpc.WithInsecure()
	}

	exporter, err := otlptrace.New(
		context.Background(),
		otlptracegrpc.NewClient(
			secureOption,
			otlptracegrpc.WithHeaders(opts.CollectorHeader),
			otlptracegrpc.WithEndpoint(opts.CollectorUrl),
		),
	)
	if err != nil {
		log.Fatal(err)
	}

	otel.SetTracerProvider(
		sdktrace.NewTracerProvider(
			sdktrace.WithSampler(sdktrace.AlwaysSample()),
			sdktrace.WithBatcher(exporter),
			sdktrace.WithResource(initResources(opts)),
		),
	)
	propagator := propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{})
	otel.SetTextMapPropagator(propagator)
	return func() {
		_ = exporter.Shutdown(context.Background())
	}
}

func InitMetrics(opts *Opts) func() {
	secureOption := otlpmetricgrpc.WithTLSCredentials(credentials.NewClientTLSFromCert(nil, ""))
	if opts.InsecureMode {
		secureOption = otlpmetricgrpc.WithInsecure()
	}

	exporter, err := otlpmetricgrpc.New(
		context.Background(),
		secureOption,
		otlpmetricgrpc.WithHeaders(opts.CollectorHeader),
		otlpmetricgrpc.WithEndpoint(opts.CollectorUrl),
	)
	if err != nil {
		log.Fatal(err)
	}

	otel.SetMeterProvider(
		sdkmetric.NewMeterProvider(
			sdkmetric.WithResource(initResources(opts)),
			sdkmetric.WithReader(
				sdkmetric.NewPeriodicReader(exporter, sdkmetric.WithInterval(30*time.Second)),
			),
		),
	)
	return func() {
		_ = exporter.Shutdown(context.Background())
	}
}
