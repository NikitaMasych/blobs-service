package data

import "blobs/internal/types"

type Blobs interface {
	New() Blobs
	Create(blob *Blob) error
	Delete(id string) error
	Select() ([]Blob, error)
	Get(id string) (*Blob, error)
}

type Blob struct {
	ID    string         `db:"id"`
	Value string         `db:"value"`
	Type  types.BlobType `db:"type"`
}
