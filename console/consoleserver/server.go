package consoleserver

import (
	"context"
	"errors"
	"net"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/zeebo/errs"
	"golang.org/x/sync/errgroup"

	"project_template/console/consoleserver/controllers"
	"project_template/dummy"
	"project_template/pkg/logger"
)

var (
	// Error is an error class that indicates internal http server error.
	Error = errs.Class("console web server error")
)

// Config contains configuration for console web server.
type Config struct {
	Address string `env:"CONSOLE_SERVER_ADDRESS" validate:"required"`
}

// Server represents console web server.
//
// architecture: Endpoint
type Server struct {
	log    logger.Logger
	config Config

	listener net.Listener
	server   http.Server

	dummyService *dummy.Service
}

// NewServer is a constructor for console web server.
func NewServer(config Config, log logger.Logger, listener net.Listener, dummyService *dummy.Service) *Server {
	server := &Server{
		log:          log,
		config:       config,
		listener:     listener,
		dummyService: dummyService,
	}

	// controllers
	dummyController := controllers.NewDummy(server.log, dummyService)

	// routes
	router := mux.NewRouter()

	// Prometheus' metrics endpoint
	router.Handle("/metrics", promhttp.Handler())

	apiRouter := router.PathPrefix("/api/v0").Subrouter()
	apiRouter.Use(server.jsonResponse)

	dummyRouter := apiRouter.PathPrefix("/dummy").Subrouter()
	dummyRouter.Use(server.withAuth)
	dummyRouter.HandleFunc("", dummyController.List).Methods(http.MethodGet)
	dummyRouter.HandleFunc("", dummyController.Create).Methods(http.MethodPost)
	dummyRouter.HandleFunc("/{id}", dummyController.Get).Methods(http.MethodGet)
	dummyRouter.HandleFunc("/{id}", dummyController.Update).Methods(http.MethodPut)
	dummyRouter.HandleFunc("/{id}", dummyController.Delete).Methods(http.MethodDelete)

	server.server = http.Server{
		Handler: router,
	}

	return server
}

// Run starts the server that host api endpoints.
func (server *Server) Run(ctx context.Context) (err error) {
	ctx, cancel := context.WithCancel(ctx)
	var group errgroup.Group
	group.Go(func() error {
		<-ctx.Done()
		return Error.Wrap(server.server.Shutdown(context.Background()))
	})
	group.Go(func() error {
		defer cancel()
		err := server.server.Serve(server.listener)
		isCancelled := errs.IsFunc(err, func(err error) bool { return errors.Is(err, context.Canceled) })
		if isCancelled || errors.Is(err, http.ErrServerClosed) {
			err = nil
		}
		return Error.Wrap(err)
	})

	return Error.Wrap(group.Wait())
}

// Close closes server and underlying listener.
func (server *Server) Close() error {
	return Error.Wrap(server.server.Close())
}

// withAuth performs initial authorization before every request.
func (server *Server) withAuth(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		/* TODO: Implement auth logic here */

		handler.ServeHTTP(w, r.Clone(ctx))
	})
}

// jsonResponse sets a response' "Content-Type" value as "application/json"
func (server *Server) jsonResponse(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		handler.ServeHTTP(w, r.Clone(r.Context()))
	})
}
