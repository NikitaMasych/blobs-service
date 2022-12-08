package data

import (
	"blobs/internal/types"
)

type Blobs interface {
	New() Blobs
	Transaction(fn func(Blobs) error) error
	Create(blob *Blob) error
	Get(id string) (*Blob, error)
	Select() ([]Blob, error)
	Delete(id string) error
}

type Blob struct {
	ID    string         `db:"id"`
	Type  types.BlobType `db:"type"`
	Value string         `db:"value"`
}
