package platform

import (
	"context"
	"encoding/json"
	"io/fs"
	"net/http"
	"strings"
	"time"
)

func Handler(assets fs.FS, ready func(context.Context) error) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) { status(w, http.StatusOK, "ok") })
	mux.HandleFunc("GET /ready", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), time.Second)
		defer cancel()
		if err := ready(ctx); err != nil {
			status(w, http.StatusServiceUnavailable, "unavailable")
			return
		}
		status(w, http.StatusOK, "ready")
	})
	files := http.FileServer(http.FS(assets))
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			w.Header().Set("Cache-Control", "no-cache")
			files.ServeHTTP(w, r)
			return
		}
		// S0 has no client routes. Unknown APIs/assets must be real 404s, not HTML.
		if !strings.HasPrefix(r.URL.Path, "/assets/") || strings.HasSuffix(r.URL.Path, ".map") {
			http.NotFound(w, r)
			return
		}
		info, err := fs.Stat(assets, strings.TrimPrefix(r.URL.Path, "/"))
		if err != nil || info.IsDir() {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
		files.ServeHTTP(w, r)
	})
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Referrer-Policy", "no-referrer")
		w.Header().Set("Content-Security-Policy", "default-src 'self'; script-src 'self'; style-src 'self'; connect-src 'self'; img-src 'self'; object-src 'none'; base-uri 'none'; frame-ancestors 'none'")
		mux.ServeHTTP(w, r)
	})
}

func status(w http.ResponseWriter, code int, value string) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(struct {
		Status string `json:"status"`
	}{value})
}
