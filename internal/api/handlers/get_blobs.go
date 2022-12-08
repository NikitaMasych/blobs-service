package handlers

import (
	"blobs/internal/api/ctx"
	"blobs/internal/data"
	"blobs/internal/resources"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"net/http"
)

func GetBlobs(w http.ResponseWriter, r *http.Request) {
	blobs, err := ctx.BlobQ(r).Select()
	if err != nil {
		ctx.Log(r).WithError(err).Error("failed to get blobs")
		ape.RenderErr(w, problems.InternalError())
		return
	}
	response := resources.BlobListResponse{
		Data: NewBlobs(blobs),
	}
	ape.Render(w, &response)
}

func NewBlobs(blobs []data.Blob) []resources.Blob {
	result := make([]resources.Blob, len(blobs))
	for i, blob := range blobs {
		result[i] = NewBlob(&blob)
	}
	return result
}
