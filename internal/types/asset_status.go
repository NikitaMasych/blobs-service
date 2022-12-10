package types

type AssetStatus int

const (
	PendingCreation AssetStatus = iota
	Created
	PendingRemoval
	Removed
)
