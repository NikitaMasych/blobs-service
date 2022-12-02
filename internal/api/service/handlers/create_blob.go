package handlers

import (
	"blobs/internal/api/contracts"
	postgres "blobs/internal/api/database"
	"blobs/internal/api/resources"
	"blobs/internal/api/types"

	"blobs/internal/api/service/requests"
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
		Log(r).WithError(err).Warn("invalid blob type")
		ape.RenderErr(w, problems.BadRequest(validation.Errors{"/data/type": errors.New("invalid blob type")})...)
		return
	}
	err = BlobQ(r).Transaction(func(blobs contracts.Blobs) error {
		if err := blobs.Create(blob); err != nil {
			return errors.Wrap(err, "failed to create blob")
		}
		return nil
	})
	if err != nil {
		// silencing error to make request idempotent
		if errors.Cause(err) != postgres.ErrBlobsConflict {
			Log(r).WithError(err).Error("failed to save blob")
			ape.RenderErr(w, problems.InternalError())
			return
		}
	}

	response := resources.BlobResponse{
		Data: NewBlob(blob),
	}
	w.WriteHeader(201)
	ape.Render(w, &response)
}

func NewBlob(blob *types.Blob) resources.Blob {
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
