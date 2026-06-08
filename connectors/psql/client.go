package psql

import (
	"context"
	"database/sql"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	pg "github.com/lib/pq"
	"github.com/opentracing/opentracing-go"
	"github.com/qustavo/sqlhooks/v2"
)

type Client struct {
	db            *sqlx.DB
	connectionURI string
	driverName    string
	tracer        opentracing.Tracer
}

func connectPostgres(ctx context.Context, connectionStr string, databaseType Driver, tracing opentracing.Tracer) (client *Client, err error) {
	pool, err := pgxpool.New(ctx, connectionStr)
	if err != nil {
		return nil, err
	}
	var db *sqlx.DB
	var driverName Driver

	if tracing != nil {
		switch driverName {
		case Postgres:
			sql.Register(opentracing_postgres, sqlhooks.Wrap(&pg.Driver{}, NewTracingHook(tracing)))
		}
	} else {
		driverName = databaseType
	}

	pool.Config().MaxConns = int32(getIntEnv("POSTGRES_MAX_CONNS", 4))
	pool.Config().MaxConnIdleTime = getDurationEnv("POSTGRES_MAX_CONN_IDLE_TIME", 30*time.Minute)
	pool.Config().MaxConnLifetime = getDurationEnv("POSTGRES_MAX_CONN_LIFETIME", time.Hour)

	db = sqlx.NewDb(stdlib.OpenDBFromPool(pool), string(driverName))
	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}

	return &Client{
		db:            db,
		connectionURI: connectionStr,
		driverName:    string(driverName),
		tracer:        tracing,
	}, nil
}

func connect(ctx context.Context, connectionStr string, databaseType Driver, tracing opentracing.Tracer) (client *Client, err error) {
	return connectPostgres(ctx, connectionStr, databaseType, tracing)
}

func NewConnection(ctx context.Context, connectionStr string) (*Client, error) {
	return connect(ctx, connectionStr, Postgres, nil)
}

func NewConnectionWithTracing(ctx context.Context, connectionStr string, databaseType Driver, tracing opentracing.Tracer) (client *Client, err error) {
	return connect(ctx, connectionStr, databaseType, tracing)
}

func (c *Client) GetClient() *sqlx.DB {
	return c.db
}

func (c *Client) GetConnectionURI() string {
	return c.connectionURI
}

func (c *Client) SetDB(db *sqlx.DB) {
	c.db = db
}

func (c *Client) Close() error {
	return c.db.Close()
}

// getIntEnv reads an integer environment variable with a fallback default.
func getIntEnv(key string, def int) int {
	val := os.Getenv(key)
	if val == "" {
		return def
	}
	n, err := strconv.Atoi(val)
	if err != nil {
		return def
	}
	if n <= 0 {
		return def
	}
	return n
}

// getDurationEnv reads a duration environment variable (e.g. "30m", "1h", "3600s") with a fallback default.
func getDurationEnv(key string, def time.Duration) time.Duration {
	val := os.Getenv(key)
	if val == "" {
		return def
	}
	d, err := time.ParseDuration(strings.TrimSpace(val))
	if err != nil {
		return def
	}
	if d <= 0 {
		return def
	}
	return d
}
