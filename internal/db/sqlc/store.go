package db

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Store defines all functions to execute db queries and transactions
type DBStore interface {
	Querier
	GetPool() *pgxpool.Pool
}

// SQLStore provides all functions to execute SQL queries and transactions
type SQLStore struct {
	connPool *pgxpool.Pool
	*Queries
}

// ContextDB wraps the pool to dynamically switch to a Tx if found in context
type ContextDB struct {
	pool *pgxpool.Pool
}

var TxKey = "db_tx"

func (db *ContextDB) Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error) {
	if tx, ok := ctx.Value(TxKey).(pgx.Tx); ok {
		return tx.Exec(ctx, sql, arguments...)
	}
	return db.pool.Exec(ctx, sql, arguments...)
}

func (db *ContextDB) Query(ctx context.Context, sql string, arguments ...interface{}) (pgx.Rows, error) {
	if tx, ok := ctx.Value(TxKey).(pgx.Tx); ok {
		return tx.Query(ctx, sql, arguments...)
	}
	return db.pool.Query(ctx, sql, arguments...)
}

func (db *ContextDB) QueryRow(ctx context.Context, sql string, arguments ...interface{}) pgx.Row {
	if tx, ok := ctx.Value(TxKey).(pgx.Tx); ok {
		return tx.QueryRow(ctx, sql, arguments...)
	}
	return db.pool.QueryRow(ctx, sql, arguments...)
}

// NewStore creates a new store
func NewStore(connPool *pgxpool.Pool) DBStore {
	return &SQLStore{
		connPool: connPool,
		Queries:  New(&ContextDB{pool: connPool}),
	}
}

func (store *SQLStore) GetPool() *pgxpool.Pool {
	return store.connPool
}
