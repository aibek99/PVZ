package connection

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

// TxOps is
type TxOps interface {
	Rollback() error
	Prepare(ctx context.Context, query string) (*sql.Stmt, error)
	Commit() error
	DB
}

// Transaction is
type Transaction struct {
	Tx *sqlx.Tx
	db Database //nolint:unused // Ignore unused field warning for db
}

// Query is
func (tx *Transaction) Query(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return tx.Tx.QueryContext(ctx, query, args...)
}

// Get is
func (tx *Transaction) Get(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return tx.Tx.GetContext(ctx, dest, query, args...)
}

// QueryRow is
func (tx *Transaction) QueryRow(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return tx.Tx.QueryRowContext(ctx, query, args...)
}

// QueryContext is
func (tx *Transaction) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return tx.QueryContext(ctx, query, args...)
}

// Select is
func (tx *Transaction) Select(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return tx.Tx.SelectContext(ctx, dest, query, args...)
}

// Execute is
func (tx *Transaction) Execute(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return tx.Tx.ExecContext(ctx, query, args...)
}

// Rollback is
func (tx *Transaction) Rollback() error {
	return tx.Tx.Rollback()
}

// Prepare is
func (tx *Transaction) Prepare(ctx context.Context, query string) (*sql.Stmt, error) {
	s, err := tx.Tx.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	return s, nil
}

// Commit is
func (tx *Transaction) Commit() error {
	return tx.Tx.Commit()
}
