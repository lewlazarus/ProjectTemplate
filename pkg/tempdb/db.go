package tempdb

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq" // using postgres driver.
	"github.com/zeebo/errs"

	"project_template"
)

// Error indicates about internal error in tembdb processing.
var Error = errs.Class("tempdb internal error")

// Execer is for executing sql.
type Execer interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
}

// TempDatabase is a database (or something that works like an isolated database,
// such as a PostgreSQL schema) with a semi-unique name which will be cleaned up
// when closed. Mainly useful for testing purposes.
type TempDatabase struct {
	conn    *sql.DB
	ConnStr string
	Schema  string
	Driver  string
	Cleanup func(*sql.DB) error
}

// OpenUnique opens a postgres database with a temporary unique schema, which will be cleaned up when closed.
func OpenUnique(ctx context.Context, config project_template.DBConfig, schemaPrefix string) (*TempDatabase, error) {
	schemaName := schemaPrefix + "-" + CreateRandomTestingSchemaName(8)
	connStr := fmt.Sprintf("postgres://%s:%s@localhost?sslmode=disable", config.User, config.Pass)
	connStrWithSchema := ConnstrWithSchema(connStr, schemaName)

	db, err := sql.Open("postgres", connStr)
	if err == nil {
		// check that connection actually worked before trying CreateSchema, to make
		// troubleshooting (lots) easier.
		err = db.PingContext(ctx)
	}
	if err != nil {
		return nil, Error.New("failed to connect to %q with driver postgres: %w", connStrWithSchema, err)
	}

	err = CreateSchema(ctx, db, schemaName)
	if err != nil {
		return nil, errs.Combine(err, db.Close())
	}

	cleanup := func(conn *sql.DB) error {
		return DropSchema(ctx, conn, schemaName)
	}

	return &TempDatabase{
		conn:    db,
		ConnStr: connStrWithSchema,
		Schema:  schemaName,
		Driver:  "postgres",
		Cleanup: cleanup,
	}, nil
}
