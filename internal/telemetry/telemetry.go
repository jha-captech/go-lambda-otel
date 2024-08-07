package telemetry

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktracer "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

type Telemeter struct {
	Logger *slog.Logger
	*sdkmetric.MeterProvider
	*sdktracer.TracerProvider
	propagation.TextMapPropagator
}

func NewTelemeter(ctx context.Context, appName string) (Telemeter, func(context.Context) error, error) {
	var (
		telemeter     Telemeter
		shutdownFuncs []func(context.Context) error
		err           error
	)

	shutdown := func(ctx context.Context) error {
		var err error
		for _, fn := range shutdownFuncs {
			err = errors.Join(err, fn(ctx))
		}
		shutdownFuncs = nil

		return err
	}

	handleErr := func(inErr error) {
		err = errors.Join(inErr, shutdown(ctx))
	}

	// propagator
	propagator := propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
	telemeter.TextMapPropagator = propagator

	// tracer
	traceExporter, err := otlptracegrpc.New(
		ctx,
		otlptracegrpc.WithEndpointURL("http://localhost:4317"),
	)
	if err != nil {
		handleErr(fmt.Errorf("[in telemetry.NewTelemeter] jaeger.New: %w", err))
		return Telemeter{}, shutdown, err
	}
	tracerProvider := sdktracer.NewTracerProvider(
		sdktracer.WithBatcher(traceExporter),
		sdktracer.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(appName),
		)),
	)
	shutdownFuncs = append(shutdownFuncs, tracerProvider.Shutdown)
	telemeter.TracerProvider = tracerProvider

	// meter
	metricExporter, err := otlpmetricgrpc.New(
		ctx,
		otlpmetricgrpc.WithEndpoint("localhost:4317"),
	)
	if err != nil {
		handleErr(fmt.Errorf("[in telemetry.NewTelemeter] jaeger.New: %w", err))
		return Telemeter{}, shutdown, nil
	}
	meterProvider := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(metricExporter)),
		sdkmetric.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(appName),
		)),
	)
	telemeter.MeterProvider = meterProvider

	return telemeter, shutdown, nil
}
