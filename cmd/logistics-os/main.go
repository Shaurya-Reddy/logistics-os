package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Shaurya-Reddy/logistics-os/internal/platform"
	"github.com/Shaurya-Reddy/logistics-os/web"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))
	if err := run(); err != nil {
		// Connection errors may contain credentials; never print driver diagnostics.
		slog.Error("application stopped", "error", err.Error())
		os.Exit(1)
	}
}

func run() error {
	if len(os.Args) != 2 {
		return fmt.Errorf("usage: logistics-os serve|migrate|probe")
	}
	command := os.Args[1]
	if command != "serve" && command != "migrate" && command != "probe" {
		return fmt.Errorf("unknown command")
	}
	if command == "probe" {
		client := http.Client{Timeout: 2 * time.Second}
		response, err := client.Get("http://127.0.0.1:8080/ready")
		if err != nil {
			return fmt.Errorf("readiness probe failed")
		}
		defer response.Body.Close()
		if response.StatusCode != http.StatusOK {
			return fmt.Errorf("application is not ready")
		}
		return nil
	}
	url := os.Getenv("DATABASE_URL")
	if command == "migrate" {
		url = os.Getenv("MIGRATION_DATABASE_URL")
	}
	if url == "" {
		return fmt.Errorf("database URL is required for %s", command)
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	pool, err := platform.OpenPool(ctx, url)
	if err != nil {
		return err
	}
	defer pool.Close()
	if command == "migrate" {
		deadline, cancel := context.WithTimeout(ctx, 30*time.Second)
		defer cancel()
		if err := platform.Migrate(deadline, pool); err != nil {
			return fmt.Errorf("migration failed; serving was not started; check database access and migration history")
		}
		slog.Info("migrations complete")
		return nil
	}
	address := os.Getenv("HTTP_ADDR")
	if address == "" {
		address = ":8080"
	}
	if _, _, err := net.SplitHostPort(address); err != nil {
		return fmt.Errorf("invalid HTTP_ADDR")
	}
	server := &http.Server{
		Addr:              address,
		Handler:           platform.Handler(web.Assets(), func(ctx context.Context) error { return platform.Ready(ctx, pool) }),
		ReadHeaderTimeout: 5 * time.Second, ReadTimeout: 10 * time.Second,
		WriteTimeout: 10 * time.Second, IdleTimeout: 60 * time.Second, MaxHeaderBytes: 16 << 10,
	}
	failures := make(chan error, 1)
	go func() { failures <- server.ListenAndServe() }()
	slog.Info("server started", "address", address)
	select {
	case err := <-failures:
		if !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("HTTP listener failed")
		}
	case <-ctx.Done():
		deadline, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(deadline); err != nil {
			_ = server.Close()
			return fmt.Errorf("graceful shutdown timed out")
		}
	}
	return nil
}
