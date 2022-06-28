package dummy

import (
	"context"
	"github.com/google/uuid"
	"time"

	"github.com/zeebo/errs"
)

// ErrDummy indicates that there was an error in the service.
var ErrDummy = errs.Class("dummy service error")

// Service is handling users related logic.
//
// architecture: Service.
type Service struct {
	dummy DB
}

// NewService is a constructor for users service.
func NewService(dummy DB) *Service {
	return &Service{
		dummy: dummy,
	}
}

// Get returns dummy item from DB.
func (service *Service) Get(ctx context.Context, id uuid.UUID) (Dummy, error) {
	user, err := service.dummy.Get(ctx, id)
	return user, ErrDummy.Wrap(err)
}

// List returns all dummy entities from DB.
func (service *Service) List(ctx context.Context) ([]Dummy, error) {
	users, err := service.dummy.List(ctx)
	return users, ErrDummy.Wrap(err)
}

// Create creates a new dummy item.
func (service *Service) Create(ctx context.Context, title string, status Status) (Dummy, error) {
	dummy := Dummy{
		ID:        uuid.New(),
		Title:     title,
		Status:    status,
		CreatedAt: time.Now(),
	}

	err := service.dummy.Create(ctx, dummy)
	if err != nil {
		return Dummy{}, ErrDummy.Wrap(err)
	}

	return dummy, nil
}

// Update updates a dummy item data.
func (service *Service) Update(ctx context.Context, id uuid.UUID, title string, status Status) error {
	err := service.dummy.Update(ctx, id, title, status)
	return ErrDummy.Wrap(err)
}

// Delete deletes a dummy item.
func (service *Service) Delete(ctx context.Context, id uuid.UUID) error {
	err := service.dummy.Delete(ctx, id)
	return ErrDummy.Wrap(err)
}
