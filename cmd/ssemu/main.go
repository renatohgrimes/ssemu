package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path"
	"ssemu"
	"ssemu/internal/app"
	"ssemu/internal/app/auth"
	"ssemu/internal/app/chat"
	"ssemu/internal/app/game"
	"ssemu/internal/app/relay"
	"ssemu/internal/database"
	"ssemu/internal/metrics"
	"ssemu/internal/network"
	"ssemu/internal/resources"
	"strconv"
	"syscall"

	"github.com/joho/godotenv"
)

func main() {
	ctx, stopContext := signal.NotifyContext(context.Background(),
		os.Interrupt, os.Kill,
		syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)
	defer stopContext()

	var err error

	if err = setupEnv(); err != nil {
		slog.Error("failed to load environment", "err", err)
		os.Exit(1)
	}

	var logfile *os.File

	if logfile, err = setupLogs(); err != nil {
		slog.Error("failed to setup logs", "err", err)
		os.Exit(1)
	}
	defer logfile.Close()

	slog.Info("starting...", "version", app.Version)

	if err = resources.Load("client"); err != nil {
		slog.Error("failed to load resources", "err", err)
		os.Exit(1)
	}

	if err = setupDb(ctx); err != nil {
		slog.Error("failed to setup database", "err", err)
		os.Exit(1)
	}
	defer database.Close(ctx)

	var authSrv, gameSrv, relaySrv, chatSrv network.Server

	if authSrv, err = createServer(ctx, app.AuthServerSettings); err != nil {
		slog.Error("failed to create auth server", "err", err)
		os.Exit(1)
	}
	defer authSrv.Shutdown()

	if gameSrv, err = createServer(ctx, app.GameServerSettings); err != nil {
		slog.Error("failed to create game server", "err", err)
		os.Exit(1)
	}
	defer gameSrv.Shutdown()

	if relaySrv, err = createServer(ctx, app.RelayServerSettings); err != nil {
		slog.Error("failed to create relay server", "err", err)
		os.Exit(1)
	}
	defer relaySrv.Shutdown()

	if chatSrv, err = createServer(ctx, app.ChatServerSettings); err != nil {
		slog.Error("failed to create chat server", "err", err)
		os.Exit(1)
	}
	defer chatSrv.Shutdown()

	err = errors.Join(
		auth.Setup(authSrv, gameSrv),
		game.Setup(gameSrv, authSrv),
		relay.Setup(relaySrv, authSrv, gameSrv),
		chat.Setup(chatSrv, gameSrv),
	)
	if err != nil {
		slog.Error("failed to setup servers", "err", err)
		os.Exit(1)
	}

	if err = metrics.Start(ctx, gameSrv); err != nil {
		slog.Error("failed to start metrics", "err", err)
		os.Exit(1)
	}
	defer metrics.Shutdown(ctx)

	httpServer := http.Server{}
	defer httpServer.Shutdown(ctx)

	go authSrv.ListenAndServe()
	go gameSrv.ListenAndServe()
	go relaySrv.ListenAndServe()
	go chatSrv.ListenAndServe()
	go serveWebPages(&httpServer)

	<-ctx.Done()

	slog.Info("shutting down...")
}

func setupEnv() error {
	if err := godotenv.Load("ssemu.env"); err != nil {
		return err
	}
	if err := godotenv.Overload("ssemu.test.env"); err == nil {
		fmt.Println("environment overloaded with testing options")
	}
	return nil
}

func setupLogs() (*os.File, error) {
	wd, _ := os.Getwd()
	logpath := path.Join(wd, "logs.txt")
	logfile, err := os.OpenFile(logpath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0o755)
	if err != nil {
		return nil, err
	}
	level := slog.LevelDebug
	if v, err := strconv.Atoi(os.Getenv("LOG_LEVEL")); err == nil {
		level = slog.Level(v)
	}
	slog.SetDefault(
		slog.New(
			slog.NewJSONHandler(
				io.MultiWriter(os.Stdout, logfile),
				&slog.HandlerOptions{Level: level})))
	return logfile, nil
}

func setupDb(ctx context.Context) error {
	return errors.Join(database.Open(ctx, app.DbSettings), database.Migrate(ctx))
}

func createServer(ctx context.Context, s network.ServerSettings) (network.Server, error) {
	address := os.Getenv("EMU_ADDRESS_PUBLIC")
	if len(address) < 1 {
		return nil, errors.New("public address cannot be empty")
	}
	return network.NewTcpServer(ctx, address, s)
}

func serveWebPages(s *http.Server) {
	const endpoint = "0.0.0.0:8000"
	m := http.NewServeMux()
	m.Handle("/", http.FileServerFS(ssemu.WebPages))
	s.Handler = m
	logger := slog.Default().With(
		slog.String("server", "ssemu.http"),
		slog.String("endpoint", endpoint),
	)
	logger.Info("listening...")
	l, err := net.Listen("tcp", endpoint)
	if err != nil {
		logger.Error("failed to start webpages listener", "err", err)
		return
	}
	if err := s.Serve(l); err != nil {
		logger.Warn("server closed", "err", err)
	}
}
