package platform

import (
	"context"
	"crypto/sha256"
	"fmt"
	"io/fs"
	"time"

	"github.com/Shaurya-Reddy/logistics-os/db/migrations"
	"github.com/Shaurya-Reddy/logistics-os/internal/platform/platformdb"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func OpenPool(ctx context.Context, url string) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, fmt.Errorf("invalid database configuration")
	}
	config.MaxConns = 4
	config.MinConns = 0
	config.ConnConfig.ConnectTimeout = time.Second
	config.MaxConnIdleTime = 5 * time.Minute
	return pgxpool.NewWithConfig(ctx, config)
}

func migrationNames(files fs.FS) ([]string, error) { return fs.Glob(files, "*.sql") }

// Ready checks the complete migration history, rejecting missing, changed and newer schemas.
func Ready(ctx context.Context, db platformdb.DBTX) error {
	return compatible(ctx, db, migrations.Files)
}

func compatible(ctx context.Context, db platformdb.DBTX, files fs.FS) error {
	names, err := migrationNames(files)
	if err != nil {
		return err
	}
	rows, err := platformdb.New(db).SchemaHistory(ctx)
	if err != nil {
		return err
	}
	if len(rows) != len(names) {
		return fmt.Errorf("incompatible schema")
	}
	for i, name := range names {
		body, err := fs.ReadFile(files, name)
		if err != nil {
			return err
		}
		if rows[i].Version != name || rows[i].Checksum != fmt.Sprintf("%x", sha256.Sum256(body)) {
			return fmt.Errorf("incompatible schema")
		}
	}
	return nil
}

func Migrate(ctx context.Context, pool *pgxpool.Pool) error {
	return migrate(ctx, pool, migrations.Files)
}

// All pending DDL and history records commit together under a transaction advisory lock.
// Only the migration role may call this; HTTP serving never applies migrations.
func migrate(ctx context.Context, pool *pgxpool.Pool, files fs.FS) error {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())
	if _, err = tx.Exec(ctx, "SELECT pg_advisory_xact_lock(732149082)"); err != nil {
		return err
	}
	var exists bool
	if err = tx.QueryRow(ctx, "SELECT to_regclass('public.schema_migrations') IS NOT NULL").Scan(&exists); err != nil {
		return err
	}
	var rows []platformdb.SchemaHistoryRow
	if exists {
		rows, err = platformdb.New(tx).SchemaHistory(ctx)
		if err != nil {
			return err
		}
	}
	names, err := migrationNames(files)
	if err != nil {
		return err
	}
	if len(rows) > len(names) {
		return fmt.Errorf("database is newer than this application")
	}
	for i, name := range names {
		body, err := fs.ReadFile(files, name)
		if err != nil {
			return err
		}
		checksum := fmt.Sprintf("%x", sha256.Sum256(body))
		if i < len(rows) {
			if rows[i].Version != name || rows[i].Checksum != checksum {
				return fmt.Errorf("migration history mismatch: %s", name)
			}
			continue
		}
		if _, err = tx.Exec(ctx, string(body), pgx.QueryExecModeSimpleProtocol); err != nil {
			return fmt.Errorf("apply %s: %w", name, err)
		}
		if _, err = tx.Exec(ctx, "INSERT INTO schema_migrations(version, checksum) VALUES ($1, $2)", name, checksum); err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}
