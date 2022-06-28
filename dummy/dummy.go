package dummy

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/zeebo/errs"
)

// ErrNoDummy indicated that user does not exist.
var ErrNoDummy = errs.Class("dummy does not exist")

type DB interface {
	// List returns all dummies from the database.
	List(ctx context.Context) ([]Dummy, error)

	// Get returns dummy by id from the database.
	Get(ctx context.Context, id uuid.UUID) (Dummy, error)

	// Create creates a dummy and writes to the database.
	Create(ctx context.Context, dummy Dummy) error

	// Update updates a dummy in the database.
	Update(ctx context.Context, id uuid.UUID, title string, status Status) error

	// Delete deletes a dummy in the database.
	Delete(ctx context.Context, id uuid.UUID) error
}

// Status defines the list of possible dummy statuses.
type Status int

const (
	// StatusActive indicates that dummy is active.
	StatusActive Status = 1
	// StatusInactive indicates that dummy is inactive.
	StatusInactive = 0
)

type Dummy struct {
	ID        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	Status    Status    `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
}
