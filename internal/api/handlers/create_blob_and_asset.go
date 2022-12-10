package handlers

import (
	"blobs/internal/api/ctx"
	"blobs/internal/api/handlers/auxiliary"
	"blobs/internal/api/requests"
	"blobs/internal/data"
	"blobs/internal/database"
	"blobs/internal/resources"
	"blobs/internal/types"
	"errors"
	"gitlab.com/distributed_lab/ape"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation"
	"gitlab.com/distributed_lab/ape/problems"
)

func CreateBlobAndAsset(w http.ResponseWriter, r *http.Request) {
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
	if err := database.NewBlobsQ(ctx.DB(r)).Create(blob); err != nil {
		ctx.Log(r).WithError(err).Error("failed to save the blob")
		ape.RenderErr(w, problems.InternalError())
		return
	}
	response := resources.BlobResponse{
		Data: auxiliary.NewBlob(blob),
	}
	w.WriteHeader(http.StatusCreated)
	ape.Render(w, &response)
	createAsset(r)
}

func createAsset(r *http.Request) {
	asset := data.Asset{
		AssetCode: types.Ordinary,
		Creator:   types.ApiAssetCreator,
		Status:    types.PendingCreation,
	}
	if err := database.NewAssetsQ(ctx.DB(r)).Create(asset); err != nil {
		ctx.Log(r).WithError(err).Error("failed to create pending creation asset")
	}
}
