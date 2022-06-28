package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/zeebo/errs"

	"project_template/dummy"
	"project_template/pkg/logger"
)

var (
	// ErrDummy is an internal error type for dummy controller.
	ErrDummy = errs.Class("dummy controller error")
)

// Dummy is a mvc controller that handles all dummy related methods.
type Dummy struct {
	log logger.Logger

	dummy *dummy.Service
}

// NewDummy is a constructor for dummy controller.
func NewDummy(log logger.Logger, dummy *dummy.Service) *Dummy {
	dummyController := &Dummy{
		log:   log,
		dummy: dummy,
	}

	return dummyController
}

func (controller *Dummy) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	result, err := controller.dummy.List(ctx)
	if err != nil {
		controller.log.Error("could not get list of dummy", ErrDummy.Wrap(err))
		controller.serveError(w, http.StatusInternalServerError, ErrDummy.Wrap(err))
		return
	}

	if result == nil {
		result = make([]dummy.Dummy, 0, 0)
	}

	if err = json.NewEncoder(w).Encode(result); err != nil {
		controller.log.Error("failed to write json response", ErrDummy.Wrap(err))
		return
	}
}

func (controller *Dummy) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	type request struct {
		Title  string       `json:"title"`
		Status dummy.Status `json:"status"`
	}

	var (
		err error
		req request
	)

	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		controller.serveError(w, http.StatusBadRequest, ErrDummy.Wrap(err))
		return
	}

	result, err := controller.dummy.Create(ctx, req.Title, req.Status)
	if err != nil {
		controller.log.Error(fmt.Sprint("could not create dummy"), ErrDummy.Wrap(err))
		controller.serveError(w, http.StatusInternalServerError, ErrDummy.Wrap(err))
		return
	}

	if err = json.NewEncoder(w).Encode(result); err != nil {
		controller.log.Error("failed to write json response", ErrDummy.Wrap(err))
		return
	}
}

func (controller *Dummy) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)

	id, err := uuid.Parse(vars["id"])
	if err != nil {
		controller.serveError(w, http.StatusBadRequest, ErrDummy.Wrap(err))
		return
	}

	result, err := controller.dummy.Get(ctx, id)

	if err != nil {
		controller.log.Error("could not get dummy", ErrDummy.Wrap(err))
		switch {
		case dummy.ErrNoDummy.Has(err):
			controller.serveError(w, http.StatusNotFound, ErrDummy.Wrap(err))
		default:
			controller.serveError(w, http.StatusInternalServerError, ErrDummy.Wrap(err))
		}
		return
	}

	if err = json.NewEncoder(w).Encode(result); err != nil {
		controller.log.Error("failed to write json response", ErrDummy.Wrap(err))
		return
	}
}

func (controller *Dummy) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)

	id, err := uuid.Parse(vars["id"])
	if err != nil {
		controller.serveError(w, http.StatusBadRequest, ErrDummy.Wrap(err))
		return
	}

	type request struct {
		Title  string       `json:"title"`
		Status dummy.Status `json:"status"`
	}

	var req request
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		controller.serveError(w, http.StatusBadRequest, ErrDummy.Wrap(err))
		return
	}

	err = controller.dummy.Update(ctx, id, req.Title, req.Status)
	if err != nil {
		controller.log.Error("could not update dummy", ErrDummy.Wrap(err))
		switch {
		case dummy.ErrNoDummy.Has(err):
			controller.serveError(w, http.StatusNotFound, ErrDummy.Wrap(err))
		default:
			controller.serveError(w, http.StatusInternalServerError, ErrDummy.Wrap(err))
		}
		return
	}
}

func (controller *Dummy) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)

	id, err := uuid.Parse(vars["id"])
	if err != nil {
		controller.serveError(w, http.StatusBadRequest, ErrDummy.Wrap(err))
		return
	}

	err = controller.dummy.Delete(ctx, id)
	if err != nil {
		controller.log.Error(fmt.Sprint("could not delete dummy"), ErrDummy.Wrap(err))
		controller.serveError(w, http.StatusInternalServerError, ErrDummy.Wrap(err))
		return
	}
}

// serveError replies to request with specific code and error.
func (controller *Dummy) serveError(w http.ResponseWriter, status int, err error) {
	w.WriteHeader(status)

	var response struct {
		Error string `json:"error"`
	}

	response.Error = err.Error()

	if err = json.NewEncoder(w).Encode(response); err != nil {
		controller.log.Error("failed to write json error response", ErrDummy.Wrap(err))
	}
}
