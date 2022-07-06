package main

import (
	"regexp"
	"strings"
)

type templateData struct {
	// RootPkgName keeps name of root package
	RootPkgName string
	// PackageName represents package name in snake_case
	PackageName string
	// EntityNamePC represents entity name in PascalCase
	EntityNamePC string
	// EntityNameCC represents entity name in camelCase
	EntityNameCC string
	// EntityNameSC represents entity name in snake_case
	EntityNameSC string
}

// pascalCaseRule describes the RegExp rule for PascalCase
var pascalCaseRule = regexp.MustCompile("^[A-Z][a-z]+(?:[A-Z][a-z]+)*$")

// primeTemplate describes the template of prime entity file
const primeTemplate = `package {{.PackageName}}

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/zeebo/errs"
)

// ErrNo{{.EntityNamePC}} indicated that user does not exist.
var ErrNo{{.EntityNamePC}} = errs.Class("{{.EntityNameCC}} does not exist")

type DB interface {
	// List returns a list of {{.EntityNameCC}} items from the database.
	List(ctx context.Context) ([]{{.EntityNamePC}}, error)

	// Get returns {{.EntityNameCC}} by id from the database.
	Get(ctx context.Context, id uuid.UUID) ({{.EntityNamePC}}, error)

	// Create creates a {{.EntityNameCC}} and writes to the database.
	Create(ctx context.Context, {{.EntityNameCC}} {{.EntityNamePC}}) error

	// Update updates a {{.EntityNameCC}} in the database.
	Update(ctx context.Context, id uuid.UUID, title string, status Status) error

	// Delete deletes a {{.EntityNameCC}} in the database.
	Delete(ctx context.Context, id uuid.UUID) error
}

// Status defines the list of possible {{.EntityNameCC}} statuses.
type Status int

const (
	// StatusActive indicates that entity is active.
	StatusActive Status = 1
	// StatusInactive indicates that entity is inactive.
	StatusInactive = 0
)

type {{.EntityNamePC}} struct {
	ID        uuid.UUID ` + "`" + `json:"id"` + "`" + `
	Title     string    ` + "`" + `json:"title"` + "`" + `
	Status    Status    ` + "`" + `json:"status"` + "`" + `
	CreatedAt time.Time ` + "`" + `json:"createdAt"` + "`" + `
}
`

// serviceTemplate describes the template of service file
const serviceTemplate = `package {{.PackageName}}

import (
	"context"
	"github.com/google/uuid"
	"time"

	"github.com/zeebo/errs"
)

// Err{{.EntityNamePC}} indicates that there was an error in the service.
var Err{{.EntityNamePC}} = errs.Class("{{.EntityNameCC}} service error")

// Service is handling {{.EntityNamePC}} related logic.
//
// architecture: Service.
type Service struct {
	{{.EntityNameCC}} DB
}

// NewService is a constructor for {{.EntityNamePC}} service.
func NewService({{.EntityNameCC}} DB) *Service {
	return &Service{
		{{.EntityNameCC}}: {{.EntityNameCC}},
	}
}

// Get returns {{.EntityNamePC}} item from DB.
func (service *Service) Get(ctx context.Context, id uuid.UUID) ({{.EntityNamePC}}, error) {
	res, err := service.{{.EntityNameCC}}.Get(ctx, id)
	return res, Err{{.EntityNamePC}}.Wrap(err)
}

// List returns all {{.EntityNamePC}} entities from DB.
func (service *Service) List(ctx context.Context) ([]{{.EntityNamePC}}, error) {
	res, err := service.{{.EntityNameCC}}.List(ctx)
	return res, Err{{.EntityNamePC}}.Wrap(err)
}

// Create creates a new {{.EntityNamePC}} item.
func (service *Service) Create(ctx context.Context, title string, status Status) ({{.EntityNamePC}}, error) {
	ent := {{.EntityNamePC}}{
		ID:        uuid.New(),
		Title:     title,
		Status:    status,
		CreatedAt: time.Now(),
	}

	err := service.{{.EntityNameCC}}.Create(ctx, ent)
	if err != nil {
		return {{.EntityNamePC}}{}, Err{{.EntityNamePC}}.Wrap(err)
	}

	return ent, nil
}

// Update updates a {{.EntityNamePC}} item data.
func (service *Service) Update(ctx context.Context, id uuid.UUID, title string, status Status) error {
	err := service.{{.EntityNameCC}}.Update(ctx, id, title, status)
	return Err{{.EntityNamePC}}.Wrap(err)
}

// Delete deletes a {{.EntityNamePC}} item.
func (service *Service) Delete(ctx context.Context, id uuid.UUID) error {
	err := service.{{.EntityNameCC}}.Delete(ctx, id)
	return Err{{.EntityNamePC}}.Wrap(err)
}
`

const dbInterfaceTemplate = `type DB interface {

// Copy starts here >>>

	// {{.EntityNamePC}} provides access to {{.EntityNameCC}} db.
	{{.EntityNamePC}}() {{.PackageName}}.DB

// Copy end here <<<
`

const dbRepoTemplate = `
package database

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/zeebo/errs"

	"{{.RootPkgName}}/{{.PackageName}}"
)

// Err{{.EntityNamePC}} indicates that there was an error in the database.
var Err{{.EntityNamePC}} = errs.Class("{{.EntityNameCC}} repository error")

// usersDB provides access to users db.
//
// architecture: Database
type {{.EntityNameCC}}DB struct {
	conn *sql.DB
}

func ({{.EntityNameCC}}DB *{{.EntityNameCC}}DB) List(ctx context.Context) ([]{{.EntityNameCC}}.{{.EntityNamePC}}, error) {
	query := ` + "`" + `SELECT id, title, status, created_at FROM {{.EntityNameCC}}` + "`" + `

	rows, err := {{.EntityNameCC}}DB.conn.QueryContext(ctx, query)
	if err != nil {
		return nil, Err{{.EntityNamePC}}.Wrap(err)
	}
	defer func() {
		err = errs.Combine(err, rows.Close())
	}()

	var res []{{.EntityNameCC}}.{{.EntityNamePC}}

	for rows.Next() {
		var item {{.EntityNameCC}}.{{.EntityNamePC}}

		err = rows.Scan(&item.ID, &item.Title, &item.Status, &item.CreatedAt)
		if err != nil {
			return nil, Err{{.EntityNamePC}}.Wrap(err)
		}

		res = append(res, item)
	}

	if err = rows.Err(); err != nil {
		return nil, Err{{.EntityNamePC}}.Wrap(err)
	}

	return res, nil
}

func ({{.EntityNameCC}}DB *{{.EntityNameCC}}DB) Get(ctx context.Context, id uuid.UUID) ({{.EntityNameCC}}.{{.EntityNamePC}}, error) {
	var res {{.EntityNameCC}}.{{.EntityNamePC}}
	query := ` + "`" + `SELECT id, title, status, created_at FROM {{.EntityNameCC}} WHERE id = $1 LIMIT 1` + "`" + `

	err := {{.EntityNameCC}}DB.conn.QueryRowContext(ctx, query, id).Scan(&res.ID, &res.Title, &res.Status, &res.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return {{.EntityNameCC}}.{{.EntityNamePC}}{}, {{.EntityNameCC}}.ErrNo{{.EntityNamePC}}.Wrap(err)
		}
	}

	return res, nil
}

func ({{.EntityNameCC}}DB *{{.EntityNameCC}}DB) Create(ctx context.Context, d {{.EntityNameCC}}.{{.EntityNamePC}}) error {
	query := ` + "`" + `INSERT INTO {{.EntityNameCC}}(id, title, status, created_at) VALUES ($1, $2, $3, $4)` + "`" + `

	_, err := {{.EntityNameCC}}DB.conn.ExecContext(ctx, query, d.ID, d.Title, d.Status, d.CreatedAt)
	return Err{{.EntityNamePC}}.Wrap(err)
}

func ({{.EntityNameCC}}DB *{{.EntityNameCC}}DB) Update(ctx context.Context, id uuid.UUID, title string, status {{.EntityNameCC}}.Status) error {
	query := "UPDATE {{.EntityNameCC}} SET title = $1, status = $2 WHERE id = $3"

	result, err := {{.EntityNameCC}}DB.conn.ExecContext(ctx, query, title, status, id)
	if err != nil {
		return Err{{.EntityNamePC}}.Wrap(err)
	}

	rowNum, err := result.RowsAffected()
	if rowNum == 0 {
		return {{.EntityNameCC}}.ErrNo{{.EntityNamePC}}.New("{{.EntityNameCC}} does not exist")
	}

	return Err{{.EntityNamePC}}.Wrap(err)

}

func ({{.EntityNameCC}}DB *{{.EntityNameCC}}DB) Delete(ctx context.Context, id uuid.UUID) error {
	query := ` + "`" + `DELETE FROM {{.EntityNameCC}} WHERE id = $1` + "`" + `

	_, err := {{.EntityNameCC}}DB.conn.ExecContext(ctx, query, id)
	return Err{{.EntityNamePC}}.Wrap(err)
}

`

// testsTemplate describes the template of file with test
const testsTemplate = `package {{.PackageName}}_test

import (
	"context"
	"testing"
	"time"
	
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"{{.RootPkgName}}"
	"{{.RootPkgName}}/database/dbtesting"
	"{{.RootPkgName}}/{{.PackageName}}"
)

func Test{{.EntityNamePC}}(t *testing.T) {

	{{.EntityNameCC}}1 := {{.PackageName}}.{{.EntityNamePC}}{
		ID:        uuid.New(),
		Title:     "val1",
		Status:    {{.PackageName}}.StatusActive,
		CreatedAt: time.Now(),
	}

	dbtesting.Run(t, func(ctx context.Context, t *testing.T, db {{.RootPkgName}}.DB) {

		{{.EntityNameCC}}Repo := db.{{.EntityNamePC}}()

		t.Run("list", func(t *testing.T) {
			_, err := {{.EntityNameCC}}Repo.List(ctx)
			require.NoError(t, err)
		})

		t.Run("create", func(t *testing.T) {
			err := {{.EntityNameCC}}Repo.Create(ctx, {{.EntityNameSC}}1)
			require.NoError(t, err)
		})

		t.Run("update", func(t *testing.T) {
			err := {{.EntityNameCC}}Repo.Update(ctx, {{.EntityNameCC}}1.ID, {{.EntityNameCC}}1.Title, {{.EntityNameCC}}1.Status)
			require.NoError(t, err)
		})

		t.Run("get", func(t *testing.T) {
			res, err := {{.EntityNameCC}}Repo.Get(ctx, {{.EntityNameCC}}1.ID)
			require.NoError(t, err)
			require.Equal(t, res.ID, {{.EntityNameCC}}1.ID)
			require.Equal(t, res.Title, {{.EntityNameCC}}1.Title)
			require.Equal(t, res.Status, {{.EntityNameCC}}1.Status)
		})
	})
}
`

const controllerTemplate = `package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/zeebo/errs"

	"{{.RootPkgName}}/{{.PackageName}}"
	"{{.RootPkgName}}/pkg/logger"
)

var (
	// Err{{.EntityNamePC}} is an internal error type for {{.EntityNameCC}} controller.
	Err{{.EntityNamePC}} = errs.Class("{{.EntityNameCC}} controller error")
)

// {{.EntityNamePC}} is a mvc controller that handles all {{.EntityNameCC}} related methods.
type {{.EntityNamePC}} struct {
	log logger.Logger

	{{.EntityNameCC}} *{{.PackageName}}.Service
}

// New{{.EntityNamePC}} is a constructor for {{.EntityNameCC}} controller.
func New{{.EntityNamePC}}(log logger.Logger, {{.EntityNameCC}} *{{.PackageName}}.Service) *{{.EntityNamePC}} {
	{{.EntityNameCC}}Controller := &{{.EntityNamePC}}{
		log:   log,
		{{.EntityNameCC}}: {{.EntityNameCC}},
	}

	return {{.EntityNameCC}}Controller
}

func (controller *{{.EntityNamePC}}) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	result, err := controller.{{.EntityNameCC}}.List(ctx)
	if err != nil {
		controller.log.Error("could not get list of {{.EntityNameCC}}", Err{{.EntityNamePC}}.Wrap(err))
		controller.serveError(w, http.StatusInternalServerError, Err{{.EntityNamePC}}.Wrap(err))
		return
	}

	if result == nil {
		result = make([]{{.EntityNameCC}}.{{.EntityNamePC}}, 0, 0)
	}

	if err = json.NewEncoder(w).Encode(result); err != nil {
		controller.log.Error("failed to write json response", Err{{.EntityNamePC}}.Wrap(err))
		return
	}
}

func (controller *{{.EntityNamePC}}) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	type request struct {
		Title  string       ` + "`" + `json:"title"` + "`" + `
		Status {{.EntityNameCC}}.Status ` + "`" + `json:"status"` + "`" + `
	}

	var (
		err error
		req request
	)

	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		controller.serveError(w, http.StatusBadRequest, Err{{.EntityNamePC}}.Wrap(err))
		return
	}

	result, err := controller.{{.EntityNameCC}}.Create(ctx, req.Title, req.Status)
	if err != nil {
		controller.log.Error(fmt.Sprint("could not create {{.EntityNameCC}}"), Err{{.EntityNamePC}}.Wrap(err))
		controller.serveError(w, http.StatusInternalServerError, Err{{.EntityNamePC}}.Wrap(err))
		return
	}

	if err = json.NewEncoder(w).Encode(result); err != nil {
		controller.log.Error("failed to write json response", Err{{.EntityNamePC}}.Wrap(err))
		return
	}
}

func (controller *{{.EntityNamePC}}) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)

	id, err := uuid.Parse(vars["id"])
	if err != nil {
		controller.serveError(w, http.StatusBadRequest, Err{{.EntityNamePC}}.Wrap(err))
		return
	}

	result, err := controller.{{.EntityNameCC}}.Get(ctx, id)

	if err != nil {
		controller.log.Error("could not get {{.EntityNameCC}}", Err{{.EntityNamePC}}.Wrap(err))
		switch {
		case {{.EntityNameCC}}.ErrNo{{.EntityNamePC}}.Has(err):
			controller.serveError(w, http.StatusNotFound, Err{{.EntityNamePC}}.Wrap(err))
		default:
			controller.serveError(w, http.StatusInternalServerError, Err{{.EntityNamePC}}.Wrap(err))
		}
		return
	}

	if err = json.NewEncoder(w).Encode(result); err != nil {
		controller.log.Error("failed to write json response", Err{{.EntityNamePC}}.Wrap(err))
		return
	}
}

func (controller *{{.EntityNamePC}}) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)

	id, err := uuid.Parse(vars["id"])
	if err != nil {
		controller.serveError(w, http.StatusBadRequest, Err{{.EntityNamePC}}.Wrap(err))
		return
	}

	type request struct {
		Title  string       ` + "`" + `json:"title"` + "`" + `
		Status {{.EntityNameCC}}.Status ` + "`" + `json:"status"` + "`" + `
	}

	var req request
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		controller.serveError(w, http.StatusBadRequest, Err{{.EntityNamePC}}.Wrap(err))
		return
	}

	err = controller.{{.EntityNameCC}}.Update(ctx, id, req.Title, req.Status)
	if err != nil {
		controller.log.Error("could not update {{.EntityNameCC}}", Err{{.EntityNamePC}}.Wrap(err))
		switch {
		case {{.EntityNameCC}}.ErrNo{{.EntityNamePC}}.Has(err):
			controller.serveError(w, http.StatusNotFound, Err{{.EntityNamePC}}.Wrap(err))
		default:
			controller.serveError(w, http.StatusInternalServerError, Err{{.EntityNamePC}}.Wrap(err))
		}
		return
	}
}

func (controller *{{.EntityNamePC}}) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)

	id, err := uuid.Parse(vars["id"])
	if err != nil {
		controller.serveError(w, http.StatusBadRequest, Err{{.EntityNamePC}}.Wrap(err))
		return
	}

	err = controller.{{.EntityNameCC}}.Delete(ctx, id)
	if err != nil {
		controller.log.Error(fmt.Sprint("could not delete {{.EntityNameCC}}"), Err{{.EntityNamePC}}.Wrap(err))
		controller.serveError(w, http.StatusInternalServerError, Err{{.EntityNamePC}}.Wrap(err))
		return
	}
}

// serveError replies to request with specific code and error.
func (controller *{{.EntityNamePC}}) serveError(w http.ResponseWriter, status int, err error) {
	w.WriteHeader(status)

	var response struct {
		Error string ` + "`" + `json:"error"` + "`" + `
	}

	response.Error = err.Error()

	if err = json.NewEncoder(w).Encode(response); err != nil {
		controller.log.Error("failed to write json error response", Err{{.EntityNamePC}}.Wrap(err))
	}
}
`

const routesTemplate = `

// Copy starts here >>>

	{{.EntityNameCC}}Router := apiRouter.PathPrefix("/{{.EntityNameCC}}").Subrouter()
	{{.EntityNameCC}}Router.Use(server.withAuth)
	{{.EntityNameCC}}Router.HandleFunc("", {{.EntityNameCC}}Controller.List).Methods(http.MethodGet)
	{{.EntityNameCC}}Router.HandleFunc("", {{.EntityNameCC}}Controller.Create).Methods(http.MethodPost)
	{{.EntityNameCC}}Router.HandleFunc("/{id}", {{.EntityNameCC}}Controller.Get).Methods(http.MethodGet)
	{{.EntityNameCC}}Router.HandleFunc("/{id}", {{.EntityNameCC}}Controller.Update).Methods(http.MethodPut)
	{{.EntityNameCC}}Router.HandleFunc("/{id}", {{.EntityNameCC}}Controller.Delete).Methods(http.MethodDelete)

// Copy end here <<<

`

// newTemplateData is the constructor for templateData
func newTemplateData(rootPkgName, entityNamePC string) *templateData {

	return &templateData{
		RootPkgName:  rootPkgName,
		PackageName:  toSnakeCase(entityNamePC),
		EntityNamePC: entityNamePC,
		EntityNameCC: toCamelCase(entityNamePC),
		EntityNameSC: toSnakeCase(entityNamePC),
	}
}

// toSnakeCase converts PascalCase string into snake_case string
func toSnakeCase(str string) string {
	matchFirstCap := regexp.MustCompile("(.)([A-Z][a-z]+)")
	matchAllCap := regexp.MustCompile("([a-z0-9])([A-Z])")

	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

// toCamelCase converts PascalCase string into camelCase string
func toCamelCase(str string) string {
	return strings.ToLower(str[0:1]) + str[1:]
}
