package postgres

import (
	"blobs/internal/api/contracts"
	"blobs/internal/api/types"
	"database/sql"
	"github.com/Masterminds/squirrel"
	"github.com/lib/pq"
	"github.com/pkg/errors"
)

const (
	blobsTable        = "blobs"
	blobsPKConstraint = "blobs_pkey"
)

var (
	ErrBlobsConflict = errors.New("blobs primary key conflict")
	blobsSelect      = squirrel.
				Select("id", "value", "type").
				From(blobsTable)
)

type Blobs struct {
	*Repo
	stmt squirrel.SelectBuilder
}

func NewBlobs(repo *Repo) *Blobs {
	return &Blobs{
		repo.Clone(), blobsSelect,
	}
}

func (q *Blobs) New() contracts.Blobs {
	return NewBlobs(q.Repo.Clone())
}

func (q *Blobs) Transaction(fn func(contracts.Blobs) error) error {
	return q.Repo.Transaction(func() error {
		return fn(q)
	})
}

func (q *Blobs) Create(blob *types.Blob) error {
	stmt := squirrel.Insert(blobsTable).SetMap(map[string]interface{}{
		"id":    blob.ID,
		"type":  blob.Type,
		"value": blob.Value,
	})

	_, err := q.Exec(stmt)
	if err != nil {
		cause := errors.Cause(err)
		pqerr, ok := cause.(*pq.Error)
		if ok {
			if pqerr.Constraint == blobsPKConstraint {
				return ErrBlobsConflict
			}
		}
	}
	return err
}

func (q *Blobs) Delete(id string) error {
	stmt := squirrel.Delete(blobsTable).Where("id = ?", id)

	_, err := q.Exec(stmt)
	if err != nil {
		cause := errors.Cause(err)
		pqerr, ok := cause.(*pq.Error)
		if ok {
			if pqerr.Constraint == blobsPKConstraint {
				return ErrBlobsConflict
			}
		}
	}
	return err
}

func (q *Blobs) Get(id string) (*types.Blob, error) {
	var result types.Blob
	stmt := q.stmt.Where("id = ?", id)

	err := q.Repo.Get(&result, stmt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &result, nil
}

func (q *Blobs) GetAll() ([]*types.Blob, error) {
	var blobs []*types.Blob
	rows, err := q.Repo.Query(blobsSelect)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	for rows.Next() {
		blob := new(types.Blob)
		if err = rows.StructScan(blob); err != nil {
			return nil, err
		}
		blobs = append(blobs, blob)
	}
	return blobs, nil

}
