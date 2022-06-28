package database

import (
	"context"
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"github.com/zeebo/errs"

	"project_template/dummy"
)

// ErrDummy indicates that there was an error in the database.
var ErrDummy = errs.Class("dummy repository error")

// usersDB provides access to users db.
//
// architecture: Database
type dummyDB struct {
	conn *sql.DB
}

func (dummyDB *dummyDB) List(ctx context.Context) ([]dummy.Dummy, error) {
	query := `SELECT id, title, status, created_at FROM dummy`

	rows, err := dummyDB.conn.QueryContext(ctx, query)
	if err != nil {
		return nil, ErrDummy.Wrap(err)
	}
	defer func() {
		err = errs.Combine(err, rows.Close())
	}()

	var result []dummy.Dummy

	for rows.Next() {
		var item dummy.Dummy

		err = rows.Scan(&item.ID, &item.Title, &item.Status, &item.CreatedAt)
		if err != nil {
			return nil, ErrDummy.Wrap(err)
		}

		result = append(result, item)
	}

	if err = rows.Err(); err != nil {
		return nil, ErrDummy.Wrap(err)
	}

	return result, nil
}

func (dummyDB *dummyDB) Get(ctx context.Context, id uuid.UUID) (dummy.Dummy, error) {
	var result dummy.Dummy
	query := `SELECT id, title, status, created_at FROM dummy WHERE id = $1 LIMIT 1`

	err := dummyDB.conn.QueryRowContext(ctx, query, id).Scan(&result.ID, &result.Title, &result.Status, &result.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return dummy.Dummy{}, dummy.ErrNoDummy.Wrap(err)
		}
	}

	return result, nil
}

func (dummyDB *dummyDB) Create(ctx context.Context, d dummy.Dummy) error {
	query := `INSERT INTO dummy(id, title, status, created_at)
	          VALUES ($1, $2, $3, $4)`

	_, err := dummyDB.conn.ExecContext(ctx, query, d.ID, d.Title, d.Status, d.CreatedAt)
	return ErrDummy.Wrap(err)
}

func (dummyDB *dummyDB) Update(ctx context.Context, id uuid.UUID, title string, status dummy.Status) error {
	query := "UPDATE dummy SET title = $1, status = $2 WHERE id = $3"

	result, err := dummyDB.conn.ExecContext(ctx, query, title, status, id)
	if err != nil {
		return ErrDummy.Wrap(err)
	}

	rowNum, err := result.RowsAffected()
	if rowNum == 0 {
		return dummy.ErrNoDummy.New("dummy does not exist")
	}

	return ErrDummy.Wrap(err)

}

func (dummyDB *dummyDB) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM dummy WHERE id = $1`

	_, err := dummyDB.conn.ExecContext(ctx, query, id)
	return ErrDummy.Wrap(err)
}
