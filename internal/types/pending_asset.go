package types

type PendingAssetStatus int

const (
	Pending PendingAssetStatus = iota
	Rejected
	Approved
)
