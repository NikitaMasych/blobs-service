package handlers

import (
	"blobs/internal/api/ctx"
	"blobs/internal/api/requests"
	"blobs/internal/resources"
	"net/http"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func GetBlob(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewGetBlobRequest(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}
	blob, err := ctx.BlobQ(r).Get(request.BlobID)
	if err != nil {
		ctx.Log(r).WithError(err).Error("failed to get blob")
		ape.RenderErr(w, problems.InternalError())
		return
	}
	if blob == nil {
		ape.RenderErr(w, problems.NotFound())
		return
	}
	response := resources.BlobResponse{
		Data: NewBlob(blob),
	}
	ape.Render(w, &response)
}
