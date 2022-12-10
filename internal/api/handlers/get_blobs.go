package handlers

import (
	"blobs/internal/api/ctx"
	"blobs/internal/api/handlers/auxiliary"
	"blobs/internal/database"
	"blobs/internal/resources"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"net/http"
)

func GetBlobs(w http.ResponseWriter, r *http.Request) {
	blobs, err := database.NewBlobsQ(ctx.DB(r)).Select()
	if err != nil {
		ctx.Log(r).WithError(err).Error("failed to get blobs")
		ape.RenderErr(w, problems.InternalError())
		return
	}
	response := resources.BlobListResponse{
		Data: auxiliary.NewBlobs(blobs),
	}
	ape.Render(w, &response)
}
