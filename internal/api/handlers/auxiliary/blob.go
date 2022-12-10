package auxiliary

import (
	"blobs/internal/data"
	"blobs/internal/resources"
)

func NewBlobs(blobs []data.Blob) []resources.Blob {
	result := make([]resources.Blob, len(blobs))
	for i, blob := range blobs {
		result[i] = NewBlob(&blob)
	}
	return result
}

func NewBlob(blob *data.Blob) resources.Blob {
	b := resources.Blob{
		Key: resources.Key{
			ID:   blob.ID,
			Type: resources.ResourceType(blob.Type.String()),
		},
		Attributes: resources.BlobAttributes{
			Value: blob.Value,
		},
	}
	return b
}
