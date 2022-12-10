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
	blobsColumns     = []string{"id", "value", "type"}
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

func (b *BlobsQ) Create(blob *data.Blob) error {
	stmt := squirrel.Insert(blobsTable).
		Columns(blobsColumns...).Values(blob.ID, blob.Value, blob.Type)

	if err := b.db.Exec(stmt); err != nil {
		cause := errors.Cause(err)
		if pqerr, ok := cause.(*pq.Error); ok {
			if pqerr.Constraint == blobsPKConstraint {
				return ErrBlobsConflict
			}
		}
		return errors.Wrap(err, "failed to create the blob")
	}
	return nil
}

func (b *BlobsQ) Delete(id string) error {
	stmt := squirrel.Delete(blobsTable).Where(squirrel.Eq{"id": id})

	if err := b.db.Exec(stmt); err != nil {
		return errors.Wrap(err, "failed to delete the blob")
	}
	return nil
}

func (b *BlobsQ) Select() ([]data.Blob, error) {
	blobs := make([]data.Blob, 0)
	if err := b.db.Select(&blobs, b.selector); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, errors.Wrap(err, "failed to select blobs")
	}
	return blobs, nil
}

func (b *BlobsQ) Get(id string) (*data.Blob, error) {
	stmt := b.selector.Where(squirrel.Eq{"id": id})
	var blob data.Blob
	if err := b.db.Get(&blob, stmt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, errors.Wrap(err, "failed to get the blob")
	}
	return &blob, nil
}
