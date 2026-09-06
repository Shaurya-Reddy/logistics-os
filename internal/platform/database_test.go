package platform

import (
	"context"
	"io/fs"
	"os"
	"sync"
	"testing"
	"testing/fstest"
	"time"

	"github.com/Shaurya-Reddy/logistics-os/db/migrations"
	"github.com/jackc/pgx/v5/pgxpool"
)

// TEST_DATABASE_URL must name a disposable empty database. No existing schema is dropped.
func TestMigrationLifecycle(t *testing.T) {
	url := os.Getenv("TEST_DATABASE_URL")
	if url == "" {
		t.Skip("set TEST_DATABASE_URL to a disposable empty PostgreSQL database")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	pool, err := pgxpool.New(ctx, url)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	var exists bool
	if err := pool.QueryRow(ctx, "SELECT to_regclass('public.schema_migrations') IS NOT NULL").Scan(&exists); err != nil {
		t.Fatal(err)
	}
	if exists {
		t.Fatal("test requires a fresh database; refusing to alter existing migration history")
	}
	if Ready(ctx, pool) == nil {
		t.Fatal("unmigrated database reported ready")
	}
	var wg sync.WaitGroup
	errors := make(chan error, 2)
	for range 2 {
		wg.Add(1)
		go func() { defer wg.Done(); errors <- Migrate(ctx, pool) }()
	}
	wg.Wait()
	close(errors)
	for err := range errors {
		if err != nil {
			t.Fatal(err)
		}
	}
	if err := Ready(ctx, pool); err != nil {
		t.Fatal(err)
	}
	if err := Migrate(ctx, pool); err != nil {
		t.Fatalf("idempotent retry: %v", err)
	}
	body, err := fs.ReadFile(migrations.Files, "0001_platform.sql")
	if err != nil {
		t.Fatal(err)
	}
	altered := fstest.MapFS{"0001_platform.sql": {Data: append(append([]byte{}, body...), '\n')}}
	if migrate(ctx, pool, altered) == nil {
		t.Fatal("changed migration accepted")
	}
	failing := fstest.MapFS{
		"0001_platform.sql": {Data: body},
		"0002_failure.sql":  {Data: []byte("CREATE TABLE rollback_probe (id integer); SELECT 1/0;")},
	}
	if migrate(ctx, pool, failing) == nil {
		t.Fatal("broken migration accepted")
	}
	if err := pool.QueryRow(ctx, "SELECT to_regclass('public.rollback_probe') IS NOT NULL").Scan(&exists); err != nil || exists {
		t.Fatalf("DDL did not roll back: %v", err)
	}
	if err := Ready(ctx, pool); err != nil {
		t.Fatalf("failure corrupted prior schema: %v", err)
	}
	if _, err := pool.Exec(ctx, "INSERT INTO schema_migrations(version, checksum) VALUES ('9999_future.sql', 'future')"); err != nil {
		t.Fatal(err)
	}
	if Ready(ctx, pool) == nil || Migrate(ctx, pool) == nil {
		t.Fatal("newer schema was accepted")
	}
}
