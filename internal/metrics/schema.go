package metrics

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "modernc.org/sqlite"
)

const schemaVersion = 1

var ErrNewerSchema = errors.New("metrics database schema is newer than this binary supports")

type store struct {
	writer *sql.DB
	reader *sql.DB
}

func openStore(ctx context.Context, databasePath string) (*store, error) {
	databasePath = strings.TrimSpace(databasePath)
	if databasePath == "" {
		return nil, errors.New("metrics database path is required")
	}
	if !filepath.IsAbs(databasePath) {
		return nil, fmt.Errorf("metrics database path must be absolute: %q", databasePath)
	}

	databasePath = filepath.Clean(databasePath)
	if err := os.MkdirAll(filepath.Dir(databasePath), 0o755); err != nil {
		return nil, fmt.Errorf("create metrics database directory: %w", err)
	}

	dsn := sqliteDSN(databasePath)
	writer, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("open metrics database writer: %w", err)
	}
	writer.SetMaxOpenConns(1)
	writer.SetMaxIdleConns(1)

	closeWriter := true
	defer func() {
		if closeWriter {
			_ = writer.Close()
		}
	}()

	if err := writer.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("connect metrics database writer: %w", err)
	}
	if err := migrate(ctx, writer); err != nil {
		return nil, err
	}

	reader, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("open metrics database reader: %w", err)
	}
	reader.SetMaxOpenConns(4)
	reader.SetMaxIdleConns(4)
	if err := reader.PingContext(ctx); err != nil {
		_ = reader.Close()
		return nil, fmt.Errorf("connect metrics database reader: %w", err)
	}

	closeWriter = false
	return &store{writer: writer, reader: reader}, nil
}

func sqliteDSN(databasePath string) string {
	values := url.Values{}
	values.Add("_pragma", "journal_mode(WAL)")
	values.Add("_pragma", "synchronous(NORMAL)")
	values.Add("_pragma", "busy_timeout(5000)")
	values.Add("_pragma", "foreign_keys(ON)")
	return (&url.URL{Scheme: "file", Path: databasePath, RawQuery: values.Encode()}).String()
}

func migrate(ctx context.Context, db *sql.DB) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin metrics schema migration: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	if _, err := tx.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version INTEGER PRIMARY KEY,
			applied_at INTEGER NOT NULL
		)`); err != nil {
		return fmt.Errorf("create metrics schema migrations table: %w", err)
	}

	var current int
	if err := tx.QueryRowContext(ctx, `SELECT COALESCE(MAX(version), 0) FROM schema_migrations`).Scan(&current); err != nil {
		return fmt.Errorf("read metrics schema version: %w", err)
	}
	if current > schemaVersion {
		return fmt.Errorf("%w: database is version %d, binary supports version %d", ErrNewerSchema, current, schemaVersion)
	}

	if current < 1 {
		if err := migrationOne(ctx, tx); err != nil {
			return err
		}
		if _, err := tx.ExecContext(ctx,
			`INSERT INTO schema_migrations(version, applied_at) VALUES (?, ?)`,
			1, time.Now().UTC().UnixNano(),
		); err != nil {
			return fmt.Errorf("record metrics schema migration 1: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit metrics schema migration: %w", err)
	}
	return nil
}

func migrationOne(ctx context.Context, tx *sql.Tx) error {
	statements := []string{
		`CREATE TABLE events (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			event_type TEXT NOT NULL,
			payload_version INTEGER NOT NULL,
			occurred_at INTEGER NOT NULL,
			created_at INTEGER NOT NULL,
			session_id TEXT,
			automation_key TEXT,
			evaluation_id TEXT,
			operation_id TEXT,
			payload_json TEXT NOT NULL
		)`,
		`CREATE INDEX events_session_id_idx ON events(session_id, id)`,
		`CREATE INDEX events_automation_key_idx ON events(automation_key, id)`,
		`CREATE INDEX events_evaluation_id_idx ON events(evaluation_id, occurred_at, id)`,
		`CREATE INDEX events_event_type_idx ON events(event_type, id)`,
		`CREATE INDEX events_occurred_at_idx ON events(occurred_at, id)`,
	}
	for _, statement := range statements {
		if _, err := tx.ExecContext(ctx, statement); err != nil {
			return fmt.Errorf("apply metrics schema migration 1: %w", err)
		}
	}
	return nil
}
