package handlers

import (
	"blobs/internal/api/ctx"
	"blobs/internal/api/requests"
	"blobs/internal/data"
	"blobs/internal/database"
	"blobs/internal/types"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"net/http"
)

func DeleteBlobAndAsset(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewGetBlobRequest(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}
	if err = database.NewBlobsQ(ctx.DB(r)).Delete(request.BlobID); err != nil {
		ctx.Log(r).WithError(err).Error("failed to delete the blob")
		ape.RenderErr(w, problems.InternalError())
		return
	}
	w.WriteHeader(http.StatusNoContent)
	removeAsset(r)
}

func removeAsset(r *http.Request) {
	asset := data.Asset{
		AssetCode: types.Ordinary,
		Creator:   types.ApiAssetCreator,
		Status:    types.PendingRemoval,
	}
	if err := database.NewAssetsQ(ctx.DB(r)).Create(asset); err != nil {
		ctx.Log(r).WithError(err).Error("failed to create pending removal asset")
	}
}
