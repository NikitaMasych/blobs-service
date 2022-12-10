package helpers

import (
	"gitlab.com/tokend/go/xdr"
)

type AssetDetails struct {
	Name  string `json:"name"`
	Owner string `json:"owner"`
}

const (
	MaxIssuanceAmount = 1000000
	OrdinaryAssetType = uint64(0)
	Decimals          = 6
	Policy            = uint32(xdr.AssetPolicyBaseAsset)
)

var ZeroTasks = uint32(0)
