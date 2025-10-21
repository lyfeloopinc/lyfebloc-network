package types

import "cosmossdk.io/collections"

const (
	ModuleName = "restaking"
	StoreKey   = ModuleName
	RouterKey  = ModuleName
)

var (
	AutoRestakeRatioKey = collections.NewPrefix("auto_restake_ratio")
)
