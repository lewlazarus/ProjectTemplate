package main

import (
	"context"
	"os"
	"project_template"
	"project_template/database"
	"project_template/pkg/config"

	"github.com/spf13/cobra"
	"github.com/zeebo/errs"

	"project_template/pkg/logger/zaplog"
)

// Error is a default error type for template_project cli.
var Error = errs.Class("template_project cli error")

// Config contains configuration for console web server.
type Config struct {
	project_template.Config
	project_template.DBConfig
}

// commands.
var (
	// template_project root cmd.
	rootCmd = &cobra.Command{
		Use:   "template_project",
		Short: "cli for interacting with template_project project",
	}

	runCmd = &cobra.Command{
		Use:         "run",
		Short:       "runs the program",
		RunE:        cmdRun,
		Annotations: map[string]string{"type": "run"},
	}
)

func init() {
	rootCmd.AddCommand(runCmd)
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func cmdRun(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	log := zaplog.NewLog()

	runCfg := Config{}
	err := config.ReadConfig(&runCfg)
	if err != nil {
		log.Error("could not read config", Error.Wrap(err))
		return Error.Wrap(err)
	}

	db, err := database.New(runCfg.DBConfig)
	if err != nil {
		log.Error("could not connect to database", Error.Wrap(err))
		return Error.Wrap(err)
	}
	defer func() {
		err = errs.Combine(err, db.Close())
	}()

	app, err := project_template.New(runCfg.Config, log, db)
	if err != nil {
		log.Error("could not start template_project service", Error.Wrap(err))
		return Error.Wrap(err)
	}

	runError := app.Run(ctx)
	closeError := app.Close()

	return Error.Wrap(errs.Combine(runError, closeError))
}
