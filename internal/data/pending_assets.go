package data

import (
	"blobs/internal/types"
)

type PendingAssets interface {
	New() PendingAssets

	Create(p PendingAsset) error
	UpdateStatus(status types.PendingAssetStatus, txId string) error
	FilterByStatus(status types.PendingAssetStatus) PendingAssets
	FilterByTxId(txId string) PendingAssets

	Get() (*PendingAsset, error)
	Select() ([]PendingAsset, error)
}

type PendingAsset struct {
	AssetCode string                   `db:"asset_code"`
	TxId      string                   `db:"tx_id"`
	Creator   string                   `db:"creator"`
	Status    types.PendingAssetStatus `db:"status"`
}
