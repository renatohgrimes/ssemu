package metrics

import (
	"context"
	"log/slog"
	"net"
	"net/http"
	"os"
	"ssemu/internal/app"
	"ssemu/internal/network"
	"ssemu/internal/telemetry"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel/metric"
)

type metricsManager struct {
	meterProvider     metric.MeterProvider
	gameServer        network.Server
	metricsHttpServer *http.Server
}

var manager *metricsManager

func Start(ctx context.Context, gameServer network.Server) error {
	if manager != nil {
		return nil
	}
	if len(os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")) == 0 {
		return nil
	}
	res := telemetry.NewResource(app.GameServerSettings.Name, app.GameServerSettings.Version)
	mp, err := telemetry.NewMeterProvider(ctx, res)
	if err != nil {
		return err
	}
	s := &http.Server{}
	go servePrometheusMetricsHttpEndpoint(s)
	meter := mp.Meter("")
	observePlayerCount(meter, gameServer)
	manager = &metricsManager{
		meterProvider:     mp,
		gameServer:        gameServer,
		metricsHttpServer: s,
	}
	return nil
}

func Shutdown(ctx context.Context) {
	if manager == nil {
		return
	}
	_ = manager.metricsHttpServer.Shutdown(ctx)
}

func servePrometheusMetricsHttpEndpoint(s *http.Server) {
	m := http.NewServeMux()
	m.Handle("/metrics", promhttp.Handler())
	s.Handler = m
	logger := slog.Default().With(
		slog.String("server", "promhttp"),
	)
	l, err := net.Listen("tcp", "0.0.0.0:9000")
	if err != nil {
		logger.Error("listen failed", "error", err)
	}
	logger.Info("listening port 9000...")
	if err := s.Serve(l); err != nil {
		logger.Warn("server closed", "result", err)
	}
}

func observePlayerCount(meter metric.Meter, server network.Server) error {
	if _, err := meter.Int64ObservableGauge(
		"player_count",
		metric.WithUnit("players"),
		metric.WithDescription("Players Online"),
		metric.WithInt64Callback(func(_ context.Context, o metric.Int64Observer) error {
			o.Observe(int64(server.GetSessionCount()))
			return nil
		}),
	); err != nil {
		return err
	}
	return nil
}
