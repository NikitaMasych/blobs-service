package handlers

import (
	"blobs/internal/api/resources"
	"blobs/internal/api/types"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"net/http"
)

func GetBlobs(w http.ResponseWriter, r *http.Request) {
	log := Log(r)
	log.Info("Request started")

	log.Info("Try to get blobs from DB")
	blobs, err := BlobQ(r).GetAll()
	if err != nil {
		Log(r).WithError(err).Error("failed to get blobs")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	log.Info("Render response")
	response := resources.BlobListResponse{
		Data: NewBlobs(blobs),
	}

	ape.Render(w, &response)
}

func NewBlobs(blobs []*types.Blob) []resources.Blob {
	result := make([]resources.Blob, len(blobs))
	for i, blob := range blobs {
		result[i] = NewBlob(blob)
	}
	return result
}
