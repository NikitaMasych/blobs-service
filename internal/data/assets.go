package data

import (
	"blobs/internal/types"
)

type Assets interface {
	New() Assets
	Create(p Asset) error
	UpdateStatus(status types.AssetStatus, assetCode types.AssetCode) error
	FilterByStatus(status types.AssetStatus) Assets
	Get() (*Asset, error)
	Select() ([]Asset, error)
}

type Asset struct {
	AssetCode types.AssetCode    `db:"asset_code"`
	Creator   types.AssetCreator `db:"creator"`
	Status    types.AssetStatus  `db:"status"`
}
