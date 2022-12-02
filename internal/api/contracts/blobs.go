package contracts

import "blobs/internal/api/types"

type Blobs interface {
	New() Blobs
	Transaction(fn func(Blobs) error) error
	Create(blob *types.Blob) error
	Get(id string) (*types.Blob, error)
	GetAll() ([]*types.Blob, error)
}
