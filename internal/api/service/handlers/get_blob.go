package handlers

import (
	"blobs/internal/api/resources"
	"blobs/internal/api/service/requests"
	"net/http"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func GetBlob(w http.ResponseWriter, r *http.Request) {
	log := Log(r)
	log.Info("Request started")
	request, err := requests.NewGetBlobRequest(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	log.Info("Try to get blob from DB")
	blob, err := BlobQ(r).Get(request.BlobID)
	if err != nil {
		Log(r).WithError(err).Error("failed to get blob")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	if blob == nil {
		ape.RenderErr(w, problems.NotFound())
		return
	}

	log.Info("Render response")
	response := resources.BlobResponse{
		Data: NewBlob(blob),
	}

	ape.Render(w, &response)
}
