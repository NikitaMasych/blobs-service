package database

import (
	"blobs/internal/data"
	"database/sql"
	"github.com/Masterminds/squirrel"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/kit/pgdb"
)

const (
	blobsTable        = "blobs"
	blobsPKConstraint = "blobs_pkey"
)

var (
	ErrBlobsConflict = errors.New("blobs primary key conflict")
	blobsColumns     = []string{"id", "type", "value"}
)

type BlobsQ struct {
	db       *pgdb.DB
	selector squirrel.SelectBuilder
}

func NewBlobsQ(db *pgdb.DB) *BlobsQ {
	return &BlobsQ{
		db:       db.Clone(),
		selector: squirrel.Select(blobsColumns...).From(blobsTable),
	}
}

func (b *BlobsQ) New() data.Blobs {
	return NewBlobsQ(b.db.Clone())
}

func (b *BlobsQ) Transaction(fn func(data.Blobs) error) error {
	return b.db.Transaction(func() error {
		return fn(b)
	})
}

func (b *BlobsQ) Create(blob *data.Blob) error {
	stmt := squirrel.Insert(blobsTable).SetMap(map[string]interface{}{
		"id":    blob.ID,
		"type":  blob.Type,
		"value": blob.Value,
	})

	err := b.db.Exec(stmt)
	if err != nil {
		cause := errors.Cause(err)
		pqerr, ok := cause.(*pq.Error)
		if ok {
			if pqerr.Constraint == blobsPKConstraint {
				return ErrBlobsConflict
			}
		}
	}
	return errors.Wrap(err, "failed to create blob")
}

func (b *BlobsQ) Delete(id string) error {
	stmt := squirrel.Delete(blobsTable).Where(squirrel.Eq{"id": id})

	err := b.db.Exec(stmt)
	if err != nil {
		cause := errors.Cause(err)
		pqerr, ok := cause.(*pq.Error)
		if ok {
			if pqerr.Constraint == blobsPKConstraint {
				return ErrBlobsConflict
			}
		}
	}
	return errors.Wrap(err, "failed to delete blob")
}

func (b *BlobsQ) Get(id string) (*data.Blob, error) {
	var result data.Blob
	stmt := b.selector.Where(squirrel.Eq{"id": id})

	err := b.db.Get(&result, stmt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, errors.Wrap(err, "failed to get blob")
	}
	return &result, nil
}

func (b *BlobsQ) Select() ([]data.Blob, error) {
	blobs := make([]data.Blob, 0)
	err := b.db.Select(blobs, b.selector)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, errors.Wrap(err, "failed to select blobs")
	}
	return blobs, nil
}
