package connection

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"

	"Homework-1/internal/config"
)

// DBops is
var _ DBops = (*TDB)(nil)

// TDB is
type TDB struct {
	Db *sqlx.DB
}

// Get is
func (d *TDB) Get(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return d.Db.GetContext(ctx, dest, query, args...)
}

// QueryRow is
func (d *TDB) QueryRow(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return d.Db.QueryRowContext(ctx, query, args...)
}

// Query is
func (d *TDB) Query(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return d.Db.QueryContext(ctx, query, args...)
}

// Select is
func (d *TDB) Select(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return d.Db.SelectContext(ctx, dest, query, args...)
}

// Execute is
func (d *TDB) Execute(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return d.Db.ExecContext(ctx, query, args...)
}

// Begin is
func (d *TDB) Begin(ctx context.Context, opts *sql.TxOptions) (TxOps, error) {
	t, err := d.Db.BeginTxx(ctx, opts)
	if err != nil {
		return nil, err
	}
	return &Transaction{Tx: t}, nil
}

// NewTDB is
func NewTDB(_ context.Context, cfgs config.Postgres, t *testing.T) (*TDB, error) {
	db, err := sqlx.Connect("postgres", GenerateDsn(cfgs))
	require.NoError(t, err)
	err = db.Ping()
	require.NoError(t, err)

	return &TDB{Db: db}, nil
}

// DropRowByID deletes a row by ID from the specified table
func (d *TDB) DropRowByID(ctx context.Context, tableName string, id int64) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1", tableName)

	_, err := d.Db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("d.Db.ExecContext: failed to delete row from %s where id = %v: %w", tableName, id, err)
	}

	return nil
}

// TruncateTable is
func (d *TDB) TruncateTable(ctx context.Context, tableName string) error {
	query := fmt.Sprintf("DELETE FROM %s;", tableName)
	_, err := d.Db.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("d.Db.ExecContext: error deleting data from table %s: %w", tableName, err)
	}

	return nil
}

// Close is
func (d *TDB) Close() error {
	return d.Db.Close()
}
