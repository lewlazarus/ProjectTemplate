package dbtesting

import (
	"context"
	"strconv"
	"strings"
	"testing"

	"project_template"
	"project_template/database"
	"project_template/pkg/config"
	"project_template/pkg/tempdb"
)

// tempMasterDB is an implementing type that cleans up after itself when closed.
type tempMasterDB struct {
	tempDB *tempdb.TempDatabase
	project_template.DB
}

// Config defines configuration for tests.
type Config struct {
	MigrationsPath string `env:"DB_MIGRATIONS_PATH" validate:"required"`
	project_template.DBConfig
}

// Run method will establish connection with db, create tables in random schema, run tests.
func Run(t *testing.T, test func(ctx context.Context, t *testing.T, db project_template.DB)) {
	t.Run("Postgres", func(t *testing.T) {
		ctx := context.Background()

		runCfg := Config{}
		err := config.ReadConfig(&runCfg)
		if err != nil {
			t.Fatal(err)
		}

		db, err := CreateMasterDB(ctx, t.Name(), "Test", 0, runCfg)
		if err != nil {
			t.Fatal(err)
		}
		defer func() {
			err := db.Close()
			if err != nil {
				t.Fatal(err)
			}
		}()

		err = db.ExecuteMigrations(ctx, runCfg.MigrationsPath, true)
		if err != nil {
			t.Fatal(err)
		}

		test(ctx, t, db)
	})
}

// CreateMasterDB creates a new DB for testing.
func CreateMasterDB(ctx context.Context, name string, category string, index int, config Config) (db project_template.DB, err error) {
	schemaSuffix := tempdb.CreateRandomTestingSchemaName(6)
	schema := SchemaName(name, category, index, schemaSuffix)

	tempDB, err := tempdb.OpenUnique(ctx, config.DBConfig, schema)
	if err != nil {
		return nil, err
	}

	return CreateMasterDBOnTopOf(tempDB)
}

// SchemaName returns a properly formatted schema string.
func SchemaName(testName, category string, index int, schemaSuffix string) string {
	// postgres has a maximum schema length of 64
	// we need additional 6 bytes for the random suffix
	//    and 4 bytes for the index "/S0/""

	indexStr := strconv.Itoa(index)

	var maxTestNameLen = 64 - len(category) - len(indexStr) - len(schemaSuffix) - 2
	if len(testName) > maxTestNameLen {
		testName = testName[:maxTestNameLen]
	}

	if schemaSuffix == "" {
		return strings.ToLower(testName + "/" + category + indexStr)
	}

	return strings.ToLower(testName + "/" + schemaSuffix + "/" + category + indexStr)
}

// CreateMasterDBOnTopOf creates a new DB on top of an already existing
// temporary database.
func CreateMasterDBOnTopOf(tempDB *tempdb.TempDatabase) (db project_template.DB, err error) {
	masterDB, err := database.NewByCoonStr(tempDB.ConnStr)
	return &tempMasterDB{DB: masterDB, tempDB: tempDB}, err
}
