package handlers

import (
	"blobs/internal/api/contracts"
	"blobs/internal/api/service/requests"
	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"net/http"
)

func DeleteBlob(w http.ResponseWriter, r *http.Request) {
	log := Log(r)
	log.Info("Request started")
	request, err := requests.NewGetBlobRequest(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	log.Info("Try to delete blob from DB")
	err = BlobQ(r).Transaction(func(blobs contracts.Blobs) error {
		if err = blobs.Delete(request.BlobID); err != nil {
			return errors.Wrap(err, "failed to delete blob")
		}
		return nil
	})
	if err != nil {
		Log(r).WithError(err).Error("failed to delete blob")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
