package platform

import (
	"context"
	"errors"
	"net/http/httptest"
	"strings"
	"testing"
	"testing/fstest"
)

func TestHealthAndStaticDuringDatabaseFailure(t *testing.T) {
	assets := fstest.MapFS{"index.html": {Data: []byte("<h1>Logistics OS</h1>")}, "assets/app-abc.js": {Data: []byte("export {}")}}
	handler := Handler(assets, func(context.Context) error { return errors.New("secret database diagnostics") })
	for _, tc := range []struct {
		path  string
		code  int
		body  string
		cache string
	}{
		{"/health", 200, `"ok"`, "no-store"},
		{"/ready", 503, `"unavailable"`, "no-store"},
		{"/", 200, "Logistics OS", "no-cache"},
		{"/assets/app-abc.js", 200, "export", "public, max-age=31536000, immutable"},
		{"/api/v1/missing", 404, "404", ""},
		{"/.vite/manifest.json", 404, "404", ""},
		{"/assets/missing.js", 404, "404", ""},
		{"/assets/", 404, "404", ""},
	} {
		t.Run(tc.path, func(t *testing.T) {
			response := httptest.NewRecorder()
			handler.ServeHTTP(response, httptest.NewRequest("GET", tc.path, nil))
			if response.Code != tc.code || !strings.Contains(response.Body.String(), tc.body) {
				t.Fatalf("unexpected response: %d %s", response.Code, response.Body.String())
			}
			if strings.Contains(response.Body.String(), "secret") {
				t.Fatal("leaked diagnostics")
			}
			if response.Header().Get("Cache-Control") != tc.cache {
				t.Fatal("incorrect cache policy")
			}
		})
	}
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, httptest.NewRequest("POST", "/ready", nil))
	if response.Code != 405 {
		t.Fatalf("POST status %d", response.Code)
	}
}

func TestReadySuccessAndDeadline(t *testing.T) {
	handler := Handler(fstest.MapFS{}, func(ctx context.Context) error {
		if _, ok := ctx.Deadline(); !ok {
			t.Error("readiness must have a deadline")
		}
		return nil
	})
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, httptest.NewRequest("GET", "/ready", nil))
	if response.Code != 200 || response.Body.String() != "{\"status\":\"ready\"}\n" {
		t.Fatal(response.Body.String())
	}
}
