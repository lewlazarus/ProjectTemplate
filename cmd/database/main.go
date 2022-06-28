package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"project_template"
	"project_template/pkg/config"
	"regexp"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/zeebo/errs"

	"project_template/database"
	"project_template/pkg/fileutils"
	"project_template/pkg/logger/zaplog"
)

// Error is a default error type for database cli.
var Error = errs.Class("database cli error")

// Config contains configurable values for the migration mechanism.
type Config struct {
	MigrationsPath string `env:"DB_MIGRATIONS_PATH" validate:"required"`

	project_template.DBConfig
}

// commands.
var (
	// database root cmd.
	rootCmd = &cobra.Command{
		Use:   "database",
		Short: "cli for interacting with project database",
	}

	// create database schema.
	createMigrationCmd = &cobra.Command{
		Use:         "create-migration",
		Short:       "creates a new migration",
		RunE:        cmdCreateMigration,
		Annotations: map[string]string{"type": "run"},
	}

	// execute database migrations.
	migrateCmd = &cobra.Command{
		Use:         "migrate",
		Short:       "executes migrations",
		RunE:        cmdMigrate,
		Annotations: map[string]string{"type": "run"},
	}

	runCfg Config
)

func init() {
	rootCmd.AddCommand(createMigrationCmd)
	rootCmd.AddCommand(migrateCmd)
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

// cmdCreateSchema creates schema for all tables and databases.
func cmdCreateMigration(cmd *cobra.Command, args []string) (err error) {
	log := zaplog.NewLog()

	if len(args) == 0 {
		log.Error("migration name is required", Error.New("invalid arguments"))
		return Error.New("invalid arguments")
	}

	runCfg := Config{}
	err = config.ReadConfig(&runCfg)
	if err != nil {
		log.Error("could not read config", Error.Wrap(err))
		return Error.Wrap(err)
	}

	fExtExpr := regexp.MustCompile(".sql$")
	curVer := 0

	files, err := ioutil.ReadDir(runCfg.MigrationsPath)
	if err != nil {
		log.Error("could not read migrations path", Error.Wrap(err))
		return err
	}
	for _, f := range files {
		if f.IsDir() {
			continue
		}

		fName := f.Name()
		r := fExtExpr.MatchString(fName)
		if r == false {
			continue
		}

		parts := strings.Split(fName, "_")
		ver, err := strconv.Atoi(parts[0])
		if err != nil {
			// Looks like that file name is without a numeric prefix.
			continue
		}
		if ver > curVer {
			curVer = ver
		}
	}

	migName := fmt.Sprintf("%06d_%s", curVer+1, args[0])
	fNames := [2]string{
		migName + ".up.sql",
		migName + ".down.sql",
	}
	for _, fName := range fNames {
		isExist, err := fileutils.IsFileExist(runCfg.MigrationsPath, fName)
		if err != nil {
			log.Error("could not check file existence ", Error.Wrap(err))
			return Error.Wrap(err)
		}
		if isExist {
			errMsg := fmt.Sprintf("File '%s' is already exists", fName)
			log.Error(errMsg, Error.New("file exists"))
			return Error.New("file exists")
		}
	}

	for _, fName := range fNames {
		if err := fileutils.CreateFile(runCfg.MigrationsPath, fName); err != nil {
			errMsg := fmt.Sprintf("could not create file '%s'", fName)
			log.Error(errMsg, Error.Wrap(err))
		} else {
			fmt.Printf("new file: %s\n", fName)
		}
	}

	return nil
}

// cmdMigrate executes migrations by path in database.
func cmdMigrate(cmd *cobra.Command, args []string) (err error) {
	var isUp bool
	ctx := context.Background()
	log := zaplog.NewLog()

	if len(args) == 0 {
		log.Error("at least 1 argument is required", Error.New("invalid arguments"))
		return Error.New("invalid arguments")
	}

	runCfg := Config{}
	err = config.ReadConfig(&runCfg)
	if err != nil {
		log.Error("could not read config", Error.Wrap(err))
		return Error.Wrap(err)
	}

	db, err := database.New(runCfg.DBConfig)
	if err != nil {
		log.Error("starting database error", Error.Wrap(err))
		return Error.Wrap(err)
	}
	defer func() {
		err = Error.Wrap(errs.Combine(err, db.Close()))
	}()

	switch args[0] {
	case "up":
		isUp = true
		if err = db.ExecuteMigrations(ctx, runCfg.MigrationsPath, isUp); err != nil {
			log.Error("migrations up error", Error.Wrap(err))
		}
	case "down":
		isUp = false
		if err = db.ExecuteMigrations(ctx, runCfg.MigrationsPath, isUp); err != nil {
			log.Error("migrations down error", Error.Wrap(err))
		}
	default:
		err = errs.New("invalid arguments")
		log.Error("invalid 1 argument", Error.Wrap(err))
	}

	return Error.Wrap(err)
}
