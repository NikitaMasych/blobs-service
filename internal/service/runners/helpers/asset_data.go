package helpers

import (
	"gitlab.com/tokend/go/xdr"
)

type AssetDetails struct {
	Name          string `json:"name"`
	ContractOwner string `json:"contract_owner"`
}

const (
	OrdinaryAssetType = uint64(0)
	Policy            = uint32(xdr.AssetPolicyBaseAsset)
)

var ZeroTasks = uint32(0)
