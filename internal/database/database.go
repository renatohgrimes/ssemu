package database

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"os"
	"ssemu"
	"ssemu/internal/telemetry"

	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/mattn/go-sqlite3"
	"github.com/uptrace/opentelemetry-go-extra/otelsql"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	oteltrace "go.opentelemetry.io/otel/trace"
)

type Settings struct {
	Dsn          string
	MaxOpenConns int
}

const (
	driver = "sqlite3"
	dbName = "ssemu.sqlite3"
)

var (
	db *sql.DB
	tp *trace.TracerProvider
)

func Open(ctx context.Context, s Settings) error {
	if db != nil {
		return nil
	}
	var err error
	tp, err = telemetry.NewTraceProvider(ctx, telemetry.NewResource("ssemu.database", "sqlite3"))
	if err != nil {
		return err
	}
	db, err = otelsql.Open(driver, s.Dsn,
		otelsql.WithDBName(dbName),
		otelsql.WithTracerProvider(tp),
		otelsql.WithAttributes(semconv.DBSystemSqlite),
	)
	if err != nil {
		return err
	}
	db.SetMaxOpenConns(s.MaxOpenConns)
	db.SetConnMaxIdleTime(30 * time.Second)
	db.SetConnMaxLifetime(1 * time.Hour)
	if err = db.Ping(); err != nil {
		return err
	}
	return nil
}

func Close(ctx context.Context) {
	if db == nil {
		return
	}
	_ = tp.Shutdown(ctx)
	_ = db.Close()
}

func GetConn() *sql.DB { return db }

func Migrate(ctx context.Context) error {
	dbDriver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		return err
	}
	source, err := iofs.New(ssemu.Migrations, "sql")
	if err != nil {
		return err
	}
	migrator, err := migrate.NewWithInstance("iofs", source, driver, dbDriver)
	if err != nil {
		return err
	}
	if err = migrator.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}
	if os.Getenv("EMU_TESTING_DATA") == "1" {
		if err := addTestData(ctx); err != nil {
			return err
		}
	}
	return nil
}

func addTestData(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if _, err := GetConn().ExecContext(ctx, ssemu.TestDataInsertCommand); err != nil {
		return err
	}
	slog.Warn("testing data has been added")
	return nil
}

func GetTracer() oteltrace.Tracer { return tp.Tracer("") }
