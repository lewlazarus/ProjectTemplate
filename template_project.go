package project_template

import (
	"context"
	"errors"
	"github.com/zeebo/errs"
	"golang.org/x/sync/errgroup"
	"net"

	"project_template/console/consoleserver"
	"project_template/dummy"
	"project_template/pkg/logger"
)

type DB interface {
	// Dummy provides access to dummy db.
	Dummy() dummy.DB

	// Close closes underlying db connection.
	Close() error

	// ExecuteMigrations executes migrations by path in database.
	ExecuteMigrations(ctx context.Context, migrationsPath string, isUp bool) error
}

type DBConfig struct {
	User string `env:"DB_USER" validate:"required"`
	Pass string `env:"DB_PASS" validate:"required"`
	Name string `env:"DB_NAME" validate:"required"`
}

// Config contains the global config.
type Config struct {

	// Console keeps the console server config
	Console struct {
		Server consoleserver.Config
	}
}

// TemplateProject is the representation of the project.
type TemplateProject struct {
	Log      logger.Logger
	Database DB

	// Dummy exposes dummy related logic.
	Dummy struct {
		Service *dummy.Service
	}

	// Console web server with web UI.
	Console struct {
		Listener net.Listener
		Endpoint *consoleserver.Server
	}
}

func New(config Config, logger logger.Logger, db DB) (*TemplateProject, error) {
	var err error

	app := &TemplateProject{
		Log:      logger,
		Database: db,
	}

	{ // dummy setup.
		app.Dummy.Service = dummy.NewService(db.Dummy())
	}

	{ // console setup.
		app.Console.Listener, err = net.Listen("tcp", config.Console.Server.Address)
		if err != nil {
			return nil, err
		}

		app.Console.Endpoint = consoleserver.NewServer(
			config.Console.Server,
			logger,
			app.Console.Listener,
			app.Dummy.Service,
		)
	}

	return app, nil
}

// Run runs project until it's either closed or it errors.
func (app *TemplateProject) Run(ctx context.Context) error {
	group, ctx := errgroup.WithContext(ctx)

	group.Go(func() error {
		return ignoreCancel(app.Console.Endpoint.Run(ctx))
	})

	return group.Wait()
}

// Close closes all the resources.
func (app *TemplateProject) Close() error {
	var errlist errs.Group

	errlist.Add(app.Console.Endpoint.Close())

	return errlist.Err()
}

// we ignore cancellation and stopping errors since they are expected.
func ignoreCancel(err error) error {
	if errors.Is(err, context.Canceled) {
		return nil
	}

	return err
}
