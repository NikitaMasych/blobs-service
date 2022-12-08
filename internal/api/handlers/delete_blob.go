package handlers

import (
	"blobs/internal/api/ctx"
	"blobs/internal/api/requests"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"net/http"
)

func DeleteBlob(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewGetBlobRequest(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}
	if err = ctx.BlobQ(r).Delete(request.BlobID); err != nil {
		ctx.Log(r).WithError(err).Error("failed to delete blob")
		ape.RenderErr(w, problems.InternalError())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
