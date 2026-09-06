package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/Shaurya-Reddy/logistics-os/internal/platform/database"
	"github.com/Shaurya-Reddy/logistics-os/internal/platform/migrate"
)

const version = "0.0.1-dev"

type config struct {
	addr          string
	databaseURL   string
	staticDir     string
	migrationsDir string
}

func loadConfig() config {
	return config{
		addr:          env("HTTP_ADDR", ":8080"),
		databaseURL:   os.Getenv("DATABASE_URL"),
		staticDir:     env("STATIC_DIR", "web/dist"),
		migrationsDir: env("MIGRATIONS_DIR", "db/migrations"),
	}
}

func env(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func main() {
	flag.Parse()
	command := "serve"
	if flag.NArg() > 0 {
		command = flag.Arg(0)
	}

	cfg := loadConfig()
	if cfg.databaseURL == "" {
		slog.Error("DATABASE_URL is required")
		os.Exit(2)
	}

	ctx := context.Background()
	pool, err := database.Open(ctx, cfg.databaseURL)
	if err != nil {
		slog.Error("database configuration failed", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	switch command {
	case "migrate":
		if err := migrate.Apply(ctx, pool, cfg.migrationsDir); err != nil {
			slog.Error("migration failed", "error", err)
			os.Exit(1)
		}
		slog.Info("migrations complete")
	case "serve":
		if err := serve(pool, cfg); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("server failed", "error", err)
			os.Exit(1)
		}
	default:
		slog.Error("unknown command", "command", command)
		os.Exit(2)
	}
}

func serve(pool database.Pool, cfg config) error {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, map[string]any{"status": "ok", "version": version})
	})
	mux.HandleFunc("GET /ready", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 750*time.Millisecond)
		defer cancel()
		if err := pool.Ping(ctx); err != nil {
			writeJSON(w, http.StatusServiceUnavailable, map[string]any{"status": "not_ready"})
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"status": "ready", "database": "connected", "version": version})
	})

	staticRoot, err := filepath.Abs(cfg.staticDir)
	if err != nil {
		return err
	}
	mux.Handle("/", spaHandler(staticRoot))

	server := &http.Server{
		Addr:              cfg.addr,
		Handler:           requestLog(mux),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-stop
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = server.Shutdown(ctx)
	}()

	slog.Info("logistics-os listening", "addr", cfg.addr)
	return server.ListenAndServe()
}

func spaHandler(root string) http.Handler {
	files := http.FileServer(http.Dir(root))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clean := filepath.Clean(r.URL.Path)
		path := filepath.Join(root, clean)
		if info, err := os.Stat(path); err == nil && !info.IsDir() {
			files.ServeHTTP(w, r)
			return
		}
		http.ServeFile(w, r, filepath.Join(root, "index.html"))
	})
}

func requestLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		slog.Info("request", "method", r.Method, "path", r.URL.Path, "duration_ms", time.Since(start).Milliseconds())
	})
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
