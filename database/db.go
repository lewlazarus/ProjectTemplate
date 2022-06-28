package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file" // using golang migrate source.
	_ "github.com/lib/pq"                                // using postgres driver.
	"github.com/zeebo/errs"

	"project_template"
	"project_template/dummy"
)

// ensures that database implements project_template.DB.
var _ project_template.DB = (*database)(nil)

var (
	// Error is the default project_template error class.
	Error = errs.Class("db error")
)

// database combines access to different database tables with a record
// of the db driver, db implementation, and db source URL.
//
// architecture: Master Database
type database struct {
	conn *sql.DB
}

// New returns project_template.DB postgresql implementation.
func New(config project_template.DBConfig) (project_template.DB, error) {
	connStr := fmt.Sprintf("postgres://%s:%s@localhost/%s?sslmode=disable", config.User, config.Pass, config.Name)
	conn, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, Error.Wrap(err)
	}

	return &database{conn: conn}, nil
}

// NewByCoonStr returns project_template.DB postgresql implementation.
func NewByCoonStr(connStr string) (project_template.DB, error) {
	conn, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, Error.Wrap(err)
	}

	return &database{conn: conn}, nil
}

// Dummy provides access to dummy db.
func (db *database) Dummy() dummy.DB {
	return &dummyDB{conn: db.conn}
}

// ExecuteMigrations executes migrations by path in database.
func (db *database) ExecuteMigrations(ctx context.Context, migrationsPath string, isUp bool) error {
	driver, err := postgres.WithInstance(db.conn, &postgres.Config{})
	if err != nil {
		return Error.Wrap(err)
	}

	m, err := migrate.NewWithDatabaseInstance("file://"+migrationsPath, "postgres", driver)
	if err != nil {
		return Error.Wrap(err)
	}

	if isUp {
		err = m.Up()
	} else {
		err = m.Down()
	}

	return Error.Wrap(err)
}

// Close closes underlying db connection.
func (db *database) Close() error {
	return Error.Wrap(db.conn.Close())
}
