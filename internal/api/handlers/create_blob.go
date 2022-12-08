package handlers

import (
	"blobs/internal/api/ctx"
	"blobs/internal/api/requests"
	"blobs/internal/data"
	"blobs/internal/resources"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func CreateBlob(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewCreateBlobRequest(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}
	blob, err := requests.Blob(request)
	if err != nil {
		ctx.Log(r).WithError(err).Warn("invalid blob type")
		ape.RenderErr(w, problems.BadRequest(
			validation.Errors{"/data/type": errors.New("invalid blob type")})...)
		return
	}
	if err := ctx.BlobQ(r).Create(blob); err != nil {
		ctx.Log(r).WithError(err).Error("failed to save blob")
		ape.RenderErr(w, problems.InternalError())
		return
	}
	response := resources.BlobResponse{
		Data: NewBlob(blob),
	}
	w.WriteHeader(http.StatusCreated)
	ape.Render(w, &response)
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
