package tempdb

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net/url"
	"strings"

	"project_template/pkg/postgres"
)

// CreateSchema creates a schema if it doesn't exist.
func CreateSchema(ctx context.Context, db Execer, schema string) (err error) {
	for try := 0; try < 5; try++ {
		_, err = db.ExecContext(ctx, `CREATE SCHEMA IF NOT EXISTS `+QuoteSchema(schema)+`;`)

		// Postgres `CREATE SCHEMA IF NOT EXISTS` may return "duplicate key value violates unique constraint".
		// In that case, we will retry rather than doing anything more complicated.
		//
		// See more in: https://stackoverflow.com/a/29908840/192220
		if postgres.IsConstraintError(err) {
			continue
		}
		return err
	}

	return err
}

// DropSchema drops the named schema.
func DropSchema(ctx context.Context, db Execer, schema string) error {
	_, err := db.ExecContext(ctx, `DROP SCHEMA `+QuoteSchema(schema)+` CASCADE;`)
	return err
}

// CreateRandomTestingSchemaName creates a random schema name string.
func CreateRandomTestingSchemaName(n int) string {
	data := make([]byte, n)
	_, err := rand.Read(data)
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(data)
}

// ConnstrWithSchema adds schema to a  connection string.
func ConnstrWithSchema(connStr, schema string) string {
	if strings.Contains(connStr, "?") {
		connStr += "&options="
	} else {
		connStr += "?options="
	}
	return connStr + url.QueryEscape("--search_path="+QuoteSchema(schema))
}

// QuoteSchema quotes a schema for use in an interpolated SQL string.
func QuoteSchema(ident string) string {
	s := strings.Replace(ident, string([]byte{0}), "", -1) // nolint
	return `"` + strings.Replace(s, `"`, `""`, -1) + `"`   // nolint
}
